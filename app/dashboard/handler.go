package dashboard

import (
	"bitbucket.org/zanvd/accountant/transaction_template"
	"database/sql"
	"html/template"
	"net/http"
)

type Handler struct {
	Database *sql.DB
}

func (dh Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	transactionTemplates, err := transaction_template.GetTransactionTemplates(dh.Database, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := template.Must(
		template.ParseFiles("templates/base.gohtml", "templates/dashboard/dashboard.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", transactionTemplates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
