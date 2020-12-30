package transaction_template

import (
	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/utility"
	"database/sql"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

const BaseUrl string = "/transaction-template/"

func AddHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		transactionTemplate := TransactionTemplate{}
		if categoryId, err := strconv.Atoi(r.FormValue("category")); err == nil {
			transactionTemplate.Category = category.Category{Id: categoryId}
		}
		if name := r.FormValue("name"); name != "" {
			transactionTemplate.Name = name
		}
		if position, err := strconv.Atoi(r.FormValue("position")); err == nil {
			transactionTemplate.Position = position
		}

		if err := InsertTransactionTemplate(db, transactionTemplate); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	categories, err := category.GetCategories(db)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	data := struct {
		Categories []category.Category
	}{
		Categories: categories,
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/add.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", data); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func DeleteHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err = DeleteTransactionTemplate(db, id); err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func EditHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	transactionTemplate, err := GetTransactionTemplate(db, id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	} else if r.Method == "POST" {
		if categoryId, err := strconv.Atoi(r.FormValue("category")); err == nil {
			transactionTemplate.Category = category.Category{Id: categoryId}
		}
		if name := r.FormValue("name"); name != "" {
			transactionTemplate.Name = name
		}
		if position, err := strconv.Atoi(r.FormValue("position")); err == nil {
			transactionTemplate.Position = position
		}

		if err = UpdateTransactionTemplate(db, transactionTemplate); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	categories, err := category.GetCategories(db)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	data := struct {
		TransactionTemplate TransactionTemplate
		Categories          []category.Category
	}{
		TransactionTemplate: transactionTemplate,
		Categories:          categories,
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/edit.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", data); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ListHandler(db *sql.DB, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactionTemplates, err := GetTransactionTemplates(db, false)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/index.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", transactionTemplates); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ViewHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	transactionTemplate, err := GetTransactionTemplate(db, id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/transaction_template/view.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", transactionTemplate); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
