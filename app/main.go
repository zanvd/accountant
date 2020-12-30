package main

import (
	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/dashboard"
	"bitbucket.org/zanvd/accountant/transaction"
	"bitbucket.org/zanvd/accountant/transaction_template"
	"bitbucket.org/zanvd/accountant/utility"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
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

	http.Handle("/", appHandler{Database: db, Handler: dashboard.Handler})

	http.Handle(transaction.BaseUrl, appHandler{Database: db, Handler: transaction.ListHandler})
	http.Handle(transaction.BaseUrl+"add/", appHandler{Database: db, Handler: transaction.AddHandler})
	http.Handle(transaction.BaseUrl+"delete/", appHandler{Database: db, Handler: transaction.DeleteHandler})
	http.Handle(transaction.BaseUrl+"edit/", appHandler{Database: db, Handler: transaction.EditHandler})
	http.Handle(transaction.BaseUrl+"view/", appHandler{Database: db, Handler: transaction.ViewHandler})

	http.Handle(transaction_template.BaseUrl, appHandler{Database: db, Handler: transaction_template.ListHandler})
	http.Handle(transaction_template.BaseUrl+"add/", appHandler{Database: db, Handler: transaction_template.AddHandler})
	http.Handle(
		transaction_template.BaseUrl+"delete/", appHandler{Database: db, Handler: transaction_template.DeleteHandler})
	http.Handle(
		transaction_template.BaseUrl+"edit/", appHandler{Database: db, Handler: transaction_template.EditHandler})
	http.Handle(
		transaction_template.BaseUrl+"view/", appHandler{Database: db, Handler: transaction_template.ViewHandler})

	http.Handle(category.BaseUrl, appHandler{Database: db, Handler: category.ListHandler})
	http.Handle(category.BaseUrl+"add/", appHandler{Database: db, Handler: category.AddHandler})
	http.Handle(category.BaseUrl+"delete/", appHandler{Database: db, Handler: category.DeleteHandler})
	http.Handle(category.BaseUrl+"edit/", appHandler{Database: db, Handler: category.EditHandler})
	http.Handle(category.BaseUrl+"view/", appHandler{Database: db, Handler: category.ViewHandler})

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}

type appHandler struct {
	Database *sql.DB
	Handler  func(*sql.DB, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := ah.Handler(ah.Database, w, r); err != nil {
		log.Printf("error occurred (%s): %+v", r.URL.RequestURI(), err)

		// Handle DB errors with a custom text.
		message := err.Error()
		if _, ok := err.(*mysql.MySQLError); ok {
			message = utility.GetMySQLErrorMessage(err)
		} else if err == sql.ErrConnDone || err == sql.ErrTxDone {
			message = "Something went wrong."
			status = http.StatusInternalServerError
		} else if err == sql.ErrNoRows {
			message = "The requested data hasn't been found."
			status = http.StatusNotFound
		}

		templateData := struct {
			ErrorMessage interface{}
			ErrorStatus  int
		}{
			ErrorMessage: message,
			ErrorStatus:  status,
		}
		w.WriteHeader(http.StatusInternalServerError)
		templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/system/error.gohtml"))
		if err := templates.ExecuteTemplate(w, "base.gohtml", templateData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
