package main

import (
	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/dashboard"
	"bitbucket.org/zanvd/accountant/transaction"
	"bitbucket.org/zanvd/accountant/transaction_template"
	"bitbucket.org/zanvd/accountant/utility"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
)

func main() {
	db, err := utility.InitDatabase()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.Handle("/", appHandler{Database: db, Handler: dashboard.Handler})

	http.Handle(category.BaseUrl, appHandler{Database: db, Handler: category.ListHandler})
	http.Handle(category.BaseUrl+"add/", appHandler{Database: db, Handler: category.AddHandler})
	http.Handle(category.BaseUrl+"delete/", appHandler{Database: db, Handler: category.DeleteHandler})
	http.Handle(category.BaseUrl+"edit/", appHandler{Database: db, Handler: category.EditHandler})
	http.Handle(category.BaseUrl+"view/", appHandler{Database: db, Handler: category.ViewHandler})

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
