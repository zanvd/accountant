package user

import (
	"database/sql"
	"log"
)

type User struct {
	Id       int
	Email    string
	Password string
	Username string
}

func CreateUserTable(db *sql.DB) {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id INT NOT NULL AUTO_INCREMENT,
            email VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
            username VARCHAR(255) NOT NULL,
            PRIMARY KEY (id),
			UNIQUE KEY (email),
			UNIQUE KEY (username)
        )
    `
	log.Println("Creating the users table.")
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if _, err := statement.Exec(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Users table created.")
}

func DeleteUser(db *sql.DB, id int) (err error) {
	_, err = db.Exec("DELETE FROM users WHERE id = ?;", id)
	return
}

func GetUser(db *sql.DB, id int) (u User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?;", id)
	err = row.Scan(&u.Id, &u.Email, &u.Password, &u.Username)
	return
}

func GetUserByUsername(db *sql.DB, username string) (u User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE username = ?;", username)
	err = row.Scan(&u.Id, &u.Email, &u.Password, &u.Username)
	return
}

func InsertUser(db *sql.DB, u User) (err error) {
	statement, err := db.Prepare("INSERT INTO users(email, password, username) VALUES (?, ?, ?);")
	if err != nil {
		return
	}
	_, err = statement.Exec(u.Email, u.Password, u.Username)
	return
}

func UpdateUser(db *sql.DB, u User) (err error) {
	statement, err := db.Prepare("UPDATE users SET email = ?, password = ?, username = ? WHERE id = ?;")
	if err != nil {
		return
	}
	_, err = statement.Exec(u.Email, u.Password, u.Username, u.Id)
	return
}
