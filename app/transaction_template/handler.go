package transaction_template

import (
	"net/http"
	"path"
	"strconv"

	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/framework"
)

const BaseUri string = "/transaction-template"

type TransactionTemplateHandler struct{}

func (TransactionTemplateHandler) GetHandlers() map[string]framework.Endpoint {
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

func (TransactionTemplateHandler) GetRoutes() map[string]string {
	return map[string]string{
		"transaction-template-add":    BaseUri + "/add",
		"transaction-template-delete": BaseUri + "/delete",
		"transaction-template-edit":   BaseUri + "/edit",
		"transaction-template-list":   BaseUri,
		"transaction-template-view":   BaseUri + "/view",
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

func AddHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		transactionTemplate := TransactionTemplate{
			UserId: rd.Session.Data.User.Id,
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
		if transType, err := strconv.Atoi(r.FormValue("type")); err == nil {
			transactionTemplate.TransactionType = TransactionType(transType)
		}

		if err := InsertTransactionTemplate(t.DB, transactionTemplate); err != nil {
			return framework.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["transaction-template-list"], http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	categories, err := category.GetCategories(t.DB, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}

	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Categories       []category.Category
			TransactionTypes map[string]TransactionType
		}{
			Categories:       categories,
			TransactionTypes: GetTransactionTypes(),
		},
		Name:  "transaction-template-add",
		Title: "Add Transaction Template",
	}
	return http.StatusOK, nil
}

func DeleteHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err = DeleteTransactionTemplate(t.DB, id, rd.Session.Data.User.Id); err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}

	http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["transaction-template-list"], http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func EditHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}

	transactionTemplate, err := GetTransactionTemplate(t.DB, id, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
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
		if transType, err := strconv.Atoi(r.FormValue("type")); err == nil {
			transactionTemplate.TransactionType = TransactionType(transType)
		}

		if err = UpdateTransactionTemplate(t.DB, transactionTemplate); err != nil {
			return framework.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["transaction-template-list"], http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	categories, err := category.GetCategories(t.DB, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}

	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Categories          []category.Category
			TransactionTemplate TransactionTemplate
			TransactionTypes    map[string]TransactionType
		}{
			Categories:          categories,
			TransactionTemplate: transactionTemplate,
			TransactionTypes:    GetTransactionTypes(),
		},
		Name:  "transaction-template-edit",
		Title: "Edit Transaction Template",
	}
	return http.StatusOK, nil
}

func ListHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactionTemplates, err := GetTransactionTemplates(t.DB, false, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			TransactionTemplates []TransactionTemplate
			TransactionTypes     map[string]TransactionType
		}{
			TransactionTemplates: transactionTemplates,
			TransactionTypes:     GetTransactionTypes(),
		},
		Name:  "transaction-template-list",
		Title: "Transaction Templates",
	}
	return http.StatusOK, nil
}

func ViewHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	transactionTemplate, err := GetTransactionTemplate(t.DB, id, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			TransactionTemplate TransactionTemplate
			TransactionTypes    map[string]TransactionType
		}{
			TransactionTemplate: transactionTemplate,
			TransactionTypes:    GetTransactionTypes(),
		},
		Name:  "transaction-template-view",
		Title: "View Transaction Template",
	}
	return http.StatusOK, nil
}
