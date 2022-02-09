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
	UserId      int
}

func CreateCategoryTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS categories (
			id INT NOT NULL AUTO_INCREMENT,
			color CHAR(7) DEFAULT '',
			description VARCHAR(150) DEFAULT '',
			name VARCHAR(30) NOT NULL,
			text_color CHAR(7) DEFAULT '#000000',
			user_id INT NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY fk_user_id (user_id)
				REFERENCES users (id)
				ON DELETE CASCADE
				ON UPDATE NO ACTION
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

func DeleteCategory(db *sql.DB, id int, uid int) (err error) {
	_, err = db.Exec("DELETE FROM categories WHERE id = ? AND user_id = ?;", id, uid)
	return
}

func GetCategory(db *sql.DB, id int, uid int) (c Category, err error) {
	row := db.QueryRow("SELECT * FROM categories WHERE id = ? AND user_id = ?;", id, uid)
	err = row.Scan(&c.Id, &c.Color, &c.Description, &c.Name, &c.TextColor, &c.UserId)
	return
}

func GetCategories(db *sql.DB, uid int) (cs []Category, err error) {
	rows, err := db.Query("SELECT * FROM categories WHERE user_id = ?;", uid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Category{}
		if err = rows.Scan(&c.Id, &c.Color, &c.Description, &c.Name, &c.TextColor, &c.UserId); err != nil {
			return
		}
		cs = append(cs, c)
	}
	return
}

func InsertCategory(db *sql.DB, c Category) (err error) {
	statement, err := db.Prepare(
		"INSERT INTO categories(color, description, name, text_color, user_id) VALUES (?, ?, ?, ?, ?);",
	)
	if err != nil {
		return
	}
	_, err = statement.Exec(c.Color, c.Description, c.Name, c.TextColor, c.UserId)
	return
}

func UpdateCategory(db *sql.DB, c Category) (err error) {
	statement, err := db.Prepare(
		"UPDATE categories SET color = ?, description = ?, name = ?, text_color = ? WHERE id = ? AND user_id = ?;",
	)
	if err != nil {
		return
	}
	_, err = statement.Exec(c.Color, c.Description, c.Name, c.TextColor, c.Id, c.UserId)
	return
}
