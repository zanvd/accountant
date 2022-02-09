package transaction_template

import (
	"database/sql"
	"log"

	"bitbucket.org/zanvd/accountant/category"
)

type TransactionTemplate struct {
	Id       int
	Category category.Category
	Name     string
	Position int
	UserId   int
}

func CreateTransactionTemplateTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS transaction_templates (
			id INT NOT NULL AUTO_INCREMENT,
			category_id INT DEFAULT NULL,
			name VARCHAR(30) NOT NULL,
			position INT NOT NULL DEFAULT 1,
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
	log.Println("Creating the transaction templates table.")
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if _, err := statement.Exec(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Transaction templates table created.")
}

func DeleteTransactionTemplate(db *sql.DB, id int, uid int) (err error) {
	_, err = db.Exec("DELETE FROM transaction_templates WHERE id = ? AND user_id = ?;", id, uid)
	return
}

func GetTransactionTemplate(db *sql.DB, id int, uid int) (tt TransactionTemplate, err error) {
	query := `
		SELECT tt.id, tt.name, tt.position, tt.user_id, c.id, c.color, c.name, c.text_color
		FROM transaction_templates tt
		LEFT JOIN categories c ON c.id = tt.category_id
		WHERE tt.id = ? AND tt.user_id = ?;
	`
	row := db.QueryRow(query, id, uid)
	err = row.Scan(
		&tt.Id, &tt.Name, &tt.Position, &tt.UserId,
		&tt.Category.Id, &tt.Category.Color, &tt.Category.Name, &tt.Category.TextColor,
	)
	return
}

func GetTransactionTemplates(db *sql.DB, orderByPosition bool, uid int) (tts []TransactionTemplate, err error) {
	query := `
		SELECT tt.id, tt.name, tt.position, tt.user_id, c.id, c.color, c.name, c.text_color
		FROM transaction_templates tt
		LEFT JOIN categories c ON c.id = tt.category_id
		WHERE tt.user_id = ?
	`
	if orderByPosition {
		query += " ORDER BY tt.position ASC;"
	} else {
		query += ";"
	}
	rows, err := db.Query(query, uid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		tt := TransactionTemplate{}
		if err = rows.Scan(
			&tt.Id, &tt.Name, &tt.Position, &tt.UserId,
			&tt.Category.Id, &tt.Category.Color, &tt.Category.Name, &tt.Category.TextColor,
		); err != nil {
			return
		}
		tts = append(tts, tt)
	}
	return
}

func InsertTransactionTemplate(db *sql.DB, tt TransactionTemplate) (err error) {
	statement, err := db.Prepare(
		"INSERT INTO transaction_templates(category_id, name, position, user_id) VALUES (?, ?, ?, ?);",
	)
	if err != nil {
		return
	}
	_, err = statement.Exec(tt.Category.Id, tt.Name, tt.Position, tt.UserId)
	return
}

func UpdateTransactionTemplate(db *sql.DB, tt TransactionTemplate) (err error) {
	statement, err := db.Prepare(`
		UPDATE transaction_templates
		SET category_id = ?, name = ?, position = ?
		WHERE id = ? AND user_id = ?;
	`)
	if err != nil {
		return
	}
	_, err = statement.Exec(tt.Category.Id, tt.Name, tt.Position, tt.Id, tt.UserId)
	return
}
