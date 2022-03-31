package stats

import (
	"net/http"

	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/transaction"
)

type StatsHandler struct{}

func (StatsHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		"/stats": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: Handler,
		},
	}
}

func (StatsHandler) GetRoutes() map[string]string {
	return map[string]string{
		"stats": "/stats",
	}
}

func (StatsHandler) GetTemplates() map[string]string {
	return map[string]string{
		"stats": "templates/stats/stats.gohtml",
	}
}

func Handler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, _ *http.Request) (int, error) {
	transactions, err := transaction.GetTransactions(t.DB, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}
	income, outcome := 0.0, 0.0
	for _, transaction := range transactions {
		if transaction.Amount > 0 {
			income += transaction.Amount
		} else {
			outcome += transaction.Amount
		}
	}

	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Income  float64
			Outcome float64
			Savings float64
		}{
			Income:  income,
			Outcome: outcome,
			Savings: income + outcome,
		},
		Name:  "stats",
		Title: "Stats",
	}

	return http.StatusOK, nil
}
