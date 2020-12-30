package transaction

import (
	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/utility"
	"database/sql"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

const BaseUrl string = "/transaction/"

func AddHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
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
		if categoryId, err := strconv.Atoi(r.FormValue("category")); err == nil {
			transaction.Category = category.Category{Id: categoryId}
		}
		transaction.Summary = r.FormValue("summary")

		if err := InsertTransaction(db, transaction); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	transaction := Transaction{}
	if name := r.URL.Query().Get("name"); name != "" {
		transaction.Name = name
	}
	if categoryId, err := strconv.Atoi(r.URL.Query().Get("category")); err == nil {
		transaction.Category = category.Category{Id: categoryId}
	}

	categories, err := category.GetCategories(db)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	data := struct {
		Categories  []category.Category
		Transaction Transaction
	}{
		Categories:  categories,
		Transaction: transaction,
	}
	templates := prepareTemplates([]string{"templates/base.gohtml", "templates/transaction/add.gohtml"})
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

	if err = DeleteTransaction(db, id); err != nil {
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

	transaction, err := GetTransaction(db, id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
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
		if categoryId, err := strconv.Atoi(r.FormValue("category")); err == nil {
			transaction.Category = category.Category{Id: categoryId}
		}
		transaction.Summary = r.FormValue("summary")

		if err = UpdateTransaction(db, transaction); err != nil {
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
		Transaction Transaction
		Categories  []category.Category
	}{
		Transaction: transaction,
		Categories:  categories,
	}
	templates := prepareTemplates([]string{"templates/base.gohtml", "templates/transaction/edit.gohtml"})
	if err = templates.ExecuteTemplate(w, "base.gohtml", data); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ListHandler(db *sql.DB, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactions, err := GetTransactions(db)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	templates := prepareTemplates([]string{"templates/base.gohtml", "templates/transaction/index.gohtml"})
	if err = templates.ExecuteTemplate(w, "base.gohtml", transactions); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ViewHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	transaction, err := GetTransaction(db, id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	templates := prepareTemplates([]string{"templates/base.gohtml", "templates/transaction/view.gohtml"})
	if err = templates.ExecuteTemplate(w, "base.gohtml", transaction); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func prepareTemplates(templates []string) *template.Template {
	return template.Must(template.New("base").Funcs(template.FuncMap{
		"dbToDisplayDate": func(dbDate string) string {
			return DbToDisplayDate(dbDate)
		},
		"today": func() string {
			return CurrentDateInDisplayFormat()
		},
	}).ParseFiles(templates...))
}
