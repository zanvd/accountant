package transaction

import (
	"database/sql"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

type AddHandler struct {
	Database *sql.DB
}

func (ah AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		transaction := Transaction{}
		if amount := r.FormValue("amount"); amount != "" {
			if floatAmount, err := strconv.ParseFloat(amount, 64); err == nil {
				transaction.Amount = floatAmount
			}
		}
		if transactionDate := r.FormValue("transaction-date"); transactionDate != "" {
			if dbDate := DisplayTimeToDb(transactionDate); dbDate != "" {
				transaction.TransactionDate = dbDate
			}
		}
		if name := r.FormValue("name"); name != "" {
			transaction.Name = name
		}
		transaction.Category = r.FormValue("category")
		transaction.Summary = r.FormValue("summary")

		if err := InsertTransaction(ah.Database, transaction); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	templates := prepareTemplates([]string{"templates/base.html", "templates/transaction/add.html"})
	if err := templates.ExecuteTemplate(w, "base.html", new(struct{})); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type DeleteHandler struct {
	Database *sql.DB
}

func (dh DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = DeleteTransaction(dh.Database, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type EditHandler struct {
	Database *sql.DB
}

func (eh EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	transaction, err := GetTransaction(eh.Database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if r.Method == "POST" {
		if amount := r.FormValue("amount"); amount != "" {
			if floatAmount, err := strconv.ParseFloat(amount, 64); err == nil {
				transaction.Amount = floatAmount
			}
		}
		if transactionDate := r.FormValue("transaction-date"); transactionDate != "" {
			if dbDate := DisplayTimeToDb(transactionDate); dbDate != "" {
				transaction.TransactionDate = dbDate
			}
		}
		if name := r.FormValue("name"); name != "" {
			transaction.Name = name
		}
		transaction.Category = r.FormValue("category")
		transaction.Summary = r.FormValue("summary")

		if err := UpdateTransaction(eh.Database, transaction); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	templates := prepareTemplates([]string{"templates/base.html", "templates/transaction/edit.html"})
	if err := templates.ExecuteTemplate(w, "base.html", transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ListHandler struct {
	Database *sql.DB
}

func (lh ListHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	transactions, err := GetTransactions(lh.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := prepareTemplates([]string{"templates/base.html", "templates/transaction/index.html"})
	if err := templates.ExecuteTemplate(w, "base.html", transactions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ViewHandler struct {
	Database *sql.DB
}

func (vh ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	transaction, err := GetTransaction(vh.Database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := prepareTemplates([]string{"templates/base.html", "templates/transaction/view.html"})
	if err := templates.ExecuteTemplate(w, "base.html", transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func prepareTemplates(templates []string) *template.Template {
	return template.Must(template.New("base").Funcs(template.FuncMap{
		"dbToDisplayDate": func(dbDate string) string {
			return DbToDisplayDate(dbDate)
		},
	}).ParseFiles(templates...))
}
