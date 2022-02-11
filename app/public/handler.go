package public

import (
	"net/http"

	"bitbucket.org/zanvd/accountant/framework"
)

type PublicHandler struct{}

func (PublicHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		"/": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: HomeHandler,
		},
	}
}

func (PublicHandler) GetRoutes() map[string]string {
	return map[string]string{
		"home": "/",
	}
}

func (PublicHandler) GetTemplates() map[string]string {
	return map[string]string{
		"home": "templates/public/home.gohtml",
	}
}

func HomeHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	t.TemplateOptions = framework.TemplateOptions{
		Name: "home",
	}
	return http.StatusOK, nil
}
