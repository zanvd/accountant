package transaction

import (
	"database/sql"
	"log"

	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/convert"
)

type Transaction struct {
	Id              int
	Amount          float64
	Category        category.Category
	Name            string
	Summary         string
	TransactionDate string
	UserId          int
}

func CreateTransactionsTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS transactions (
			id INT NOT NULL AUTO_INCREMENT,
			amount DOUBLE(14, 4) NOT NULL,
			category_id INT DEFAULT NULL,
			name VARCHAR(30) NOT NULL,
			summary VARCHAR(200) DEFAULT '',
			transaction_date DATE NOT NULL,
			user_id INT NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY fk_category_id (category_id)
			    REFERENCES categories (id)
			    ON DELETE RESTRICT
				ON UPDATE NO ACTION,
			FOREIGN KEY fk_user_id (user_id)
				REFERENCES users (id)
				ON DELETE CASCADE
				ON UPDATE NO ACTION
		);
	`
	log.Println("Creating the transactions table.")
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if _, err := statement.Exec(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Transactions table created.")
}

func DeleteTransaction(db *sql.DB, id int, uid int) (err error) {
	_, err = db.Exec("DELETE FROM transactions WHERE id = ? AND user_id = ?;", id, uid)
	return
}

func GetTransaction(db *sql.DB, id int, uid int) (t Transaction, err error) {
	query := `
		SELECT t.id, t.amount, t.name, t.summary, t.transaction_date, t.user_id,
		       c.id, c.color, c.description, c.name, c.text_color
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.id = ? AND t.user_id = ?;
	`
	row := db.QueryRow(query, id, uid)
	err = row.Scan(
		&t.Id, &t.Amount, &t.Name, &t.Summary, &t.TransactionDate, &t.UserId,
		&t.Category.Id, &t.Category.Color, &t.Category.Description, &t.Category.Name, &t.Category.TextColor,
	)
	return
}

func GetTransactions(db *sql.DB, uid int) (ts []Transaction, err error) {
	query := `
		SELECT t.id, t.amount, t.name, t.summary, t.transaction_date, t.user_id,
		       c.id, c.color, c.description, c.name, c.text_color
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = ?;
	`
	rows, err := db.Query(query, uid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := Transaction{}
		if err = rows.Scan(
			&t.Id, &t.Amount, &t.Name, &t.Summary, &t.TransactionDate, &t.UserId,
			&t.Category.Id, &t.Category.Color, &t.Category.Description, &t.Category.Name, &t.Category.TextColor,
		); err != nil {
			return
		}
		ts = append(ts, t)
	}
	return
}

func InsertTransaction(db *sql.DB, t Transaction) (err error) {
	if len(t.TransactionDate) == 0 {
		t.TransactionDate = convert.CurrentDateInDbFormat()
	}
	statement, err := db.Prepare(`
		INSERT INTO transactions(amount, category_id, name, summary, transaction_date, user_id)
		VALUES (?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return
	}
	_, err = statement.Exec(t.Amount, t.Category.Id, t.Name, t.Summary, t.TransactionDate, t.UserId)
	return
}

func UpdateTransaction(db *sql.DB, t Transaction) (err error) {
	statement, err := db.Prepare(`
		UPDATE transactions
		SET amount = ?, category_id = ?, name = ?, summary = ?, transaction_date = ?
		WHERE id = ? AND user_id = ?;
	`)
	if err != nil {
		return
	}
	_, err = statement.Exec(t.Amount, t.Category.Id, t.Name, t.Summary, t.TransactionDate, t.Id, t.UserId)
	return
}
