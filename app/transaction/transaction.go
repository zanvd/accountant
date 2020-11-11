package transaction

import (
	"database/sql"
	"log"
	"time"
)

const (
	dbDateFormat      = "2006-01-02"
	displayDateFormat = "02. 01. 2006"
)

type Transaction struct {
	Id              int
	Amount          float64
	Category        string
	Name            string
	Summary         string
	TransactionDate string
}

func CreateTransactionsTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS transactions (
			id INT NOT NULL AUTO_INCREMENT,
			amount DOUBLE(14, 4) NOT NULL,
			category VARCHAR(100) DEFAULT '',
			name VARCHAR(30) NOT NULL,
			summary VARCHAR(200) DEFAULT '',
			transaction_date DATE NOT NULL,
			PRIMARY KEY (id)
		);
	`
	log.Println("Creating transactions table.")
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if _, err := statement.Exec(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Transactions table created.")
}

func DbToDisplayDate(dbDate string) string {
	return changeDateFormat(displayDateFormat, dbDateFormat, dbDate)
}

func DisplayTimeToDb(displayDate string) string {
	return changeDateFormat(dbDateFormat, displayDateFormat, displayDate)
}

func DeleteTransaction(db *sql.DB, id int) (err error) {
	_, err = db.Exec("DELETE FROM transactions WHERE id = ?", id)
	return
}

func GetTransaction(db *sql.DB, id int) (transaction Transaction, err error) {
	row := db.QueryRow("SELECT * FROM transactions WHERE id = ?;", id)
	err = row.Scan(&transaction.Id, &transaction.Amount, &transaction.Category, &transaction.Name, &transaction.Summary,
		&transaction.TransactionDate)
	if err != nil {
		return
	}
	return
}

func GetTransactions(db *sql.DB) (transactions []Transaction, err error) {
	rows, err := db.Query("SELECT * FROM transactions;")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		if err = rows.Scan(&transaction.Id, &transaction.Amount, &transaction.Category, &transaction.Name,
			&transaction.Summary, &transaction.TransactionDate); err != nil {
			return
		}
		transactions = append(transactions, transaction)
	}
	return
}

func InsertTransaction(db *sql.DB, transaction Transaction) (err error) {
	if len(transaction.TransactionDate) == 0 {
		transaction.TransactionDate = time.Now().UTC().Format(dbDateFormat)
	}
	query := "INSERT INTO transactions(amount, category, name, summary, transaction_date) VALUES (?, ?, ?, ?, ?);"
	statement, err := db.Prepare(query)
	if err != nil {
		return
	}
	_, err = statement.Exec(transaction.Amount, transaction.Category, transaction.Name, transaction.Summary,
		transaction.TransactionDate)
	return
}

func UpdateTransaction(db *sql.DB, transaction Transaction) (err error) {
	query := `
		UPDATE transactions
		SET amount = ?, category = ?, name = ?, summary = ?, transaction_date = ?
		WHERE id = ?
	`
	statement, err := db.Prepare(query)
	if err != nil {
		return
	}
	_, err = statement.Exec(transaction.Amount, transaction.Category, transaction.Name, transaction.Summary,
		transaction.TransactionDate, transaction.Id)
	return
}

func changeDateFormat(destFormat string, sourceFormat string, sourceDate string) (destDate string) {
	t, err := time.Parse(sourceFormat, sourceDate)
	if err != nil {
		return
	}
	destDate = t.Format(destFormat)
	return
}
