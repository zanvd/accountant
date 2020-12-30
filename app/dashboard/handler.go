package dashboard

import (
	"bitbucket.org/zanvd/accountant/transaction_template"
	"bitbucket.org/zanvd/accountant/utility"
	"database/sql"
	"html/template"
	"net/http"
)

func Handler(db *sql.DB, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactionTemplates, err := transaction_template.GetTransactionTemplates(db, true)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	templates := template.Must(
		template.ParseFiles("templates/base.gohtml", "templates/dashboard/dashboard.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", transactionTemplates); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
