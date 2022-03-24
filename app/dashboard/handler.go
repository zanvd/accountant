package dashboard

import (
	"net/http"

	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/transaction_template"
)

type DashboardHandler struct{}

func (DashboardHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		"/dashboard": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: Handler,
		},
	}
}

func (DashboardHandler) GetRoutes() map[string]string {
	return map[string]string{
		"dashboard": "/dashboard",
	}
}

func (DashboardHandler) GetTemplates() map[string]string {
	return map[string]string{
		"dashboard": "templates/dashboard/dashboard.gohtml",
	}
}

func Handler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactionTemplates, err := transaction_template.GetTransactionTemplates(t.DB, true, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data:  transactionTemplates,
		Name:  "dashboard",
		Title: "Dashboard",
	}
	return http.StatusOK, nil
}
