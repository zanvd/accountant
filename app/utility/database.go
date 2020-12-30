package utility

import (
	"database/sql"
	"log"
	"os"
)

func InitDatabase() (*sql.DB, error) {
	dbName, err := readSecretValueFromFile(os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalln("Failed to get the DB name:", err)
	}
	dbPassword, err := readSecretValueFromFile(os.Getenv("DB_PASSWORD"))
	if err != nil {
		log.Fatalln("Failed to get the DB password:", err)
	}
	dbUser, err := readSecretValueFromFile(os.Getenv("DB_USER"))
	if err != nil {
		log.Fatalln("Failed to get the DB user:", err)
	}
	return sql.Open("mysql", dbUser+":"+dbPassword+"@(accountant_database)/"+dbName)
}
