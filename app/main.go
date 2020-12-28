package main

import (
	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/dashboard"
	"bitbucket.org/zanvd/accountant/transaction"
	"bitbucket.org/zanvd/accountant/transaction_template"
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

	category.CreateCategoryTable(db)
	transaction.CreateTransactionsTable(db)
	transaction_template.CreateTransactionTemplateTable(db)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.Handle("/", dashboard.Handler{Database: db})

	http.Handle(transaction.BaseUrl, transaction.ListHandler{Database: db})
	http.Handle(transaction.BaseUrl+"add/", transaction.AddHandler{Database: db})
	http.Handle(transaction.BaseUrl+"delete/", transaction.DeleteHandler{Database: db})
	http.Handle(transaction.BaseUrl+"edit/", transaction.EditHandler{Database: db})
	http.Handle(transaction.BaseUrl+"view/", transaction.ViewHandler{Database: db})

	http.Handle(transaction_template.BaseUrl, transaction_template.ListHandler{Database: db})
	http.Handle(transaction_template.BaseUrl+"add/", transaction_template.AddHandler{Database: db})
	http.Handle(transaction_template.BaseUrl+"delete/", transaction_template.DeleteHandler{Database: db})
	http.Handle(transaction_template.BaseUrl+"edit/", transaction_template.EditHandler{Database: db})
	http.Handle(transaction_template.BaseUrl+"view/", transaction_template.ViewHandler{Database: db})

	http.Handle(category.BaseUrl, category.ListHandler{Database: db})
	http.Handle(category.BaseUrl+"add/", category.AddHandler{Database: db})
	http.Handle(category.BaseUrl+"delete/", category.DeleteHandler{Database: db})
	http.Handle(category.BaseUrl+"edit/", category.EditHandler{Database: db})
	http.Handle(category.BaseUrl+"view/", category.ViewHandler{Database: db})

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
