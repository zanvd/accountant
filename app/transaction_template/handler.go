package transaction_template

import (
	"bitbucket.org/zanvd/accountant/category"
	"database/sql"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

const BaseUrl string = "/transaction-template/"

type AddHandler struct {
	Database *sql.DB
}

func (ah AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		transactionTemplate := TransactionTemplate{}
		if categoryId, err := strconv.Atoi(r.FormValue("category")); err == nil {
			transactionTemplate.Category = category.Category{Id: categoryId}
		}
		if name := r.FormValue("name"); name != "" {
			transactionTemplate.Name = name
		}

		if err := InsertTransactionTemplate(ah.Database, transactionTemplate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	}

	categories, err := category.GetCategories(ah.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		Categories []category.Category
	}{
		Categories: categories,
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/add.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", data); err != nil {
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

	if err = DeleteTransactionTemplate(dh.Database, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
}

type EditHandler struct {
	Database *sql.DB
}

func (eh EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	transactionTemplate, err := GetTransactionTemplate(eh.Database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if r.Method == "POST" {
		if categoryId, err := strconv.Atoi(r.FormValue("category")); err == nil {
			transactionTemplate.Category = category.Category{Id: categoryId}
		}
		if name := r.FormValue("name"); name != "" {
			transactionTemplate.Name = name
		}

		if err := UpdateTransactionTemplate(eh.Database, transactionTemplate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	}

	categories, err := category.GetCategories(eh.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		TransactionTemplate TransactionTemplate
		Categories          []category.Category
	}{
		TransactionTemplate: transactionTemplate,
		Categories:          categories,
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/edit.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ListHandler struct {
	Database *sql.DB
}

func (lh ListHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	transactionTemplates, err := GetTransactionTemplates(lh.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/index.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", transactionTemplates); err != nil {
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
	transactionTemplate, err := GetTransactionTemplate(vh.Database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/view.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", transactionTemplate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
