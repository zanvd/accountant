package transaction_template

import (
	"net/http"
	"path"
	"strconv"

	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/utility"
)

const BaseUrl string = "/transaction-template"

type TransactionTemplateHandler struct{}

func (TransactionTemplateHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		BaseUrl: {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ListHandler,
		},
		BaseUrl + "/add": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: AddHandler,
		},
		BaseUrl + "/delete/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: DeleteHandler,
		},
		BaseUrl + "/edit/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: EditHandler,
		},
		BaseUrl + "/view/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ViewHandler,
		},
	}
}

func (TransactionTemplateHandler) GetTemplates() map[string]string {
	return map[string]string{
		"transaction-template-add":  "templates/transaction_template/add.gohtml",
		"transaction-template-edit": "templates/transaction_template/edit.gohtml",
		"transaction-template-list": "templates/transaction_template/index.gohtml",
		"transaction-template-view": "templates/transaction_template/view.gohtml",
	}
}

func AddHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		transactionTemplate := TransactionTemplate{
			UserId: t.Session.Data.User.Id,
		}
		if categoryId, err := strconv.Atoi(r.FormValue("category")); err == nil {
			transactionTemplate.Category = category.Category{Id: categoryId}
		}
		if name := r.FormValue("name"); name != "" {
			transactionTemplate.Name = name
		}
		if position, err := strconv.Atoi(r.FormValue("position")); err == nil {
			transactionTemplate.Position = position
		}

		if err := InsertTransactionTemplate(t.DB, transactionTemplate); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	categories, err := category.GetCategories(t.DB, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	t.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Categories []category.Category
		}{
			Categories: categories,
		},
		Name:  "transaction-template-add",
		Title: "Add Transaction Template",
	}
	return http.StatusOK, nil
}

func DeleteHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err = DeleteTransactionTemplate(t.DB, id, t.Session.Data.User.Id); err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func EditHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	transactionTemplate, err := GetTransactionTemplate(t.DB, id, t.Session.Data.User.Id)
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

		if err = UpdateTransactionTemplate(t.DB, transactionTemplate); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	categories, err := category.GetCategories(t.DB, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	t.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			TransactionTemplate TransactionTemplate
			Categories          []category.Category
		}{
			TransactionTemplate: transactionTemplate,
			Categories:          categories,
		},
		Name:  "transaction-template-edit",
		Title: "Edit Transaction Template",
	}
	return http.StatusOK, nil
}

func ListHandler(t *framework.Tools, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactionTemplates, err := GetTransactionTemplates(t.DB, false, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	t.TemplateOptions = framework.TemplateOptions{
		Data:  transactionTemplates,
		Name:  "transaction-template-list",
		Title: "Transaction Templates",
	}
	return http.StatusOK, nil
}

func ViewHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	transactionTemplate, err := GetTransactionTemplate(t.DB, id, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	t.TemplateOptions = framework.TemplateOptions{
		Data:  transactionTemplate,
		Name:  "transaction-template-view",
		Title: "View Transaction Template",
	}
	return http.StatusOK, nil
}
