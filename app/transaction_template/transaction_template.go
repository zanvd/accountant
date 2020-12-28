package transaction_template

import (
	"bitbucket.org/zanvd/accountant/category"
	"database/sql"
	"log"
)

type TransactionTemplate struct {
	Id       int
	Category category.Category
	Name     string
	Position int
}

func CreateTransactionTemplateTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS transaction_templates (
			id INT NOT NULL AUTO_INCREMENT,
			category_id INT DEFAULT NULL,
			name VARCHAR(30) NOT NULL,
			position INT NOT NULL DEFAULT 1,
			PRIMARY KEY (id),
			FOREIGN KEY category_id_idx (category_id)
			REFERENCES categories (id)
			ON DELETE RESTRICT
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

func DeleteTransactionTemplate(db *sql.DB, id int) (err error) {
	_, err = db.Exec("DELETE FROM transaction_templates WHERE id = ?", id)
	return
}

func GetTransactionTemplate(db *sql.DB, id int) (transactionTemplate TransactionTemplate, err error) {
	query := `
		SELECT tt.id, tt.name, tt.position, c.id, c.color, c.name, c.text_color
		FROM transaction_templates tt
		LEFT JOIN categories c ON c.id = tt.category_id
		WHERE tt.id = ?;
	`
	row := db.QueryRow(query, id)
	err = row.Scan(&transactionTemplate.Id, &transactionTemplate.Name, &transactionTemplate.Position,
		&transactionTemplate.Category.Id, &transactionTemplate.Category.Color, &transactionTemplate.Category.Name,
		&transactionTemplate.Category.TextColor)
	return
}

func GetTransactionTemplates(db *sql.DB, orderByPosition bool) (transactionTemplates []TransactionTemplate, err error) {
	query := `
		SELECT tt.id, tt.name, c.id, c.color, c.name, c.text_color
		FROM transaction_templates tt
		LEFT JOIN categories c ON c.id = tt.category_id
	`
	if orderByPosition {
		query += " ORDER BY position ASC;"
	} else {
		query += ";"
	}
	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		transactionTemplate := TransactionTemplate{}
		if err = rows.Scan(&transactionTemplate.Id, &transactionTemplate.Name, &transactionTemplate.Category.Id,
			&transactionTemplate.Category.Color, &transactionTemplate.Category.Name,
			&transactionTemplate.Category.TextColor); err != nil {
			return
		}
		transactionTemplates = append(transactionTemplates, transactionTemplate)
	}
	return
}

func InsertTransactionTemplate(db *sql.DB, transactionTemplate TransactionTemplate) (err error) {
	query := "INSERT INTO transaction_templates(category_id, name, position) VALUES (?, ?, ?)"
	statement, err := db.Prepare(query)
	if err != nil {
		return
	}
	_, err = statement.Exec(transactionTemplate.Category.Id, transactionTemplate.Name, transactionTemplate.Position)
	return
}

func UpdateTransactionTemplate(db *sql.DB, transactionTemplate TransactionTemplate) (err error) {
	query := `
		UPDATE transaction_templates
		SET category_id = ?, name = ?, position = ?
		WHERE id = ?;
	`
	statement, err := db.Prepare(query)
	if err != nil {
		return
	}
	_, err = statement.Exec(
		transactionTemplate.Category.Id, transactionTemplate.Name, transactionTemplate.Position, transactionTemplate.Id)
	return
}
