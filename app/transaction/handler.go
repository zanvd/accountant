package transaction

import (
	"net/http"
	"path"
	"strconv"

	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/convert"
	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/utility"
)

const BaseUri string = "/transaction"

type TransactionHandler struct{}

func (TransactionHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		BaseUri: {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ListHandler,
		},
		BaseUri + "/add": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: AddHandler,
		},
		BaseUri + "/delete/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: DeleteHandler,
		},
		BaseUri + "/edit/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: EditHandler,
		},
		BaseUri + "/view/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ViewHandler,
		},
	}
}

func (TransactionHandler) GetRoutes() map[string]string {
	return map[string]string{
		"transaction-add":    BaseUri + "/add",
		"transaction-delete": BaseUri + "/delete",
		"transaction-edit":   BaseUri + "/edit",
		"transaction-list":   BaseUri,
		"transaction-view":   BaseUri + "/view",
	}
}

func (TransactionHandler) GetTemplates() map[string]string {
	return map[string]string{
		"transaction-add":  "templates/transaction/add.gohtml",
		"transaction-edit": "templates/transaction/edit.gohtml",
		"transaction-list": "templates/transaction/index.gohtml",
		"transaction-view": "templates/transaction/view.gohtml",
	}
}

func AddHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		transaction := Transaction{
			UserId: t.Session.Data.User.Id,
		}
		if amount := r.FormValue("amount"); amount != "" {
			if floatAmount, err := strconv.ParseFloat(amount, 64); err == nil {
				transaction.Amount = floatAmount
			}
		}
		if transactionDate := r.FormValue("transaction-date"); transactionDate != "" {
			if dbDate := convert.DisplayTimeToDb(transactionDate); dbDate != "" {
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

		if err := InsertTransaction(t.DB, transaction); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["transaction-list"], http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	transaction := Transaction{}
	if name := r.URL.Query().Get("name"); name != "" {
		transaction.Name = name
	}
	if categoryId, err := strconv.Atoi(r.URL.Query().Get("category")); err == nil {
		transaction.Category = category.Category{Id: categoryId}
	}

	categories, err := category.GetCategories(t.DB, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	t.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Categories  []category.Category
			Transaction Transaction
		}{
			Categories:  categories,
			Transaction: transaction,
		},
		Name:  "transaction-add",
		Title: "Add Transaction",
	}
	return http.StatusOK, nil
}

func DeleteHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err = DeleteTransaction(t.DB, id, t.Session.Data.User.Id); err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["transaction-list"], http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func EditHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	transaction, err := GetTransaction(t.DB, id, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	} else if r.Method == "POST" {
		if amount := r.FormValue("amount"); amount != "" {
			if floatAmount, err := strconv.ParseFloat(amount, 64); err == nil {
				transaction.Amount = floatAmount
			}
		}
		if transactionDate := r.FormValue("transaction-date"); transactionDate != "" {
			if dbDate := convert.DisplayTimeToDb(transactionDate); dbDate != "" {
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

		if err = UpdateTransaction(t.DB, transaction); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["transaction-list"], http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	categories, err := category.GetCategories(t.DB, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	t.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Transaction Transaction
			Categories  []category.Category
		}{
			Transaction: transaction,
			Categories:  categories,
		},
		Name:  "transaction-edit",
		Title: "Edit Transaction",
	}
	return http.StatusOK, nil
}

func ListHandler(t *framework.Tools, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactions, err := GetTransactions(t.DB, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	t.TemplateOptions = framework.TemplateOptions{
		Data:  transactions,
		Name:  "transaction-list",
		Title: "Transactions",
	}
	return http.StatusOK, nil
}

func ViewHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	transaction, err := GetTransaction(t.DB, id, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	t.TemplateOptions = framework.TemplateOptions{
		Data:  transaction,
		Name:  "transaction-view",
		Title: "View Transaction",
	}
	return http.StatusOK, nil
}
