package category

import (
	"database/sql"
	"log"
)

const (
	defaultColor     = "#f5f5dc"
	defaultTextColor = "#000000"
)

type Category struct {
	Id          int
	Color       string
	Description string
	Name        string
	TextColor   string
}

func CreateCategoryTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS categories (
			id INT NOT NULL AUTO_INCREMENT,
			color CHAR(7) DEFAULT '',
			description VARCHAR(150) DEFAULT '',
			name VARCHAR(30) NOT NULL,
			text_color CHAR(7) DEFAULT '#000000',
			PRIMARY KEY (id)
		);
	`
	log.Println("Creating the categories table.")
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if _, err := statement.Exec(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Categories table created.")
}

func DeleteCategory(db *sql.DB, id int) (err error) {
	_, err = db.Exec("DELETE FROM categories WHERE id = ?;", id)
	return
}

func GetCategory(db *sql.DB, id int) (category Category, err error) {
	row := db.QueryRow("SELECT * FROM categories WHERE id = ?;", id)
	err = row.Scan(&category.Id, &category.Color, &category.Description, &category.Name, &category.TextColor)
	return
}

func GetCategories(db *sql.DB) (categories []Category, err error) {
	rows, err := db.Query("SELECT  * FROM categories;")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		category := Category{}
		if err = rows.Scan(&category.Id, &category.Color, &category.Description, &category.Name, &category.TextColor); err != nil {
			return
		}
		categories = append(categories, category)
	}
	return
}

func InsertCategory(db *sql.DB, category Category) (err error) {
	statement, err := db.Prepare(
		"INSERT INTO categories(color, description, name, text_color) VALUES (?, ?, ?, ?);",
	)
	if err != nil {
		return
	}
	_, err = statement.Exec(category.Color, category.Description, category.Name, category.TextColor)
	return
}

func UpdateCategory(db *sql.DB, category Category) (err error) {
	statement, err := db.Prepare(
		"UPDATE categories SET color = ?, description = ?, name = ?, text_color = ? WHERE id = ?;",
	)
	if err != nil {
		return
	}
	_, err = statement.Exec(category.Color, category.Description, category.Name, category.TextColor, category.Id)
	return
}
