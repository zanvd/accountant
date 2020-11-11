package main

import (
	"bitbucket.org/zanvd/accountant/transaction"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	db, err := initDatabase()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	transaction.CreateTransactionsTable(db)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.Handle("/", transaction.ListHandler{Database: db})
	http.Handle("/transaction/add/", transaction.AddHandler{Database: db})
	http.Handle("/transaction/delete/", transaction.DeleteHandler{Database: db})
	http.Handle("/transaction/edit/", transaction.EditHandler{Database: db})
	http.Handle("/transaction/view/", transaction.ViewHandler{Database: db})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}

// Deprecated. Creates a file for the SQLite DB.
func createDbFile(removeExisting bool) string {
	fileName := "accountant-sqlite.db"
	if removeExisting {
		_ = os.Remove(fileName)
	}
	log.Println("Creating DB file.")
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err := file.Close(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(fileName, "created.")
	return fileName
}

func initDatabase() (*sql.DB, error) {
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

func readSecretValueFromFile(filePath string) (secret string, err error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("path to secret is not a file: %s", filePath)
	}
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	return strings.TrimSpace(string(buffer)), nil
}
