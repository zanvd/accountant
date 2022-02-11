package user

import (
	"net/http"

	"bitbucket.org/zanvd/accountant/framework"
)

type UserHandler struct{}

func (UserHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		"/profile": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ProfileHandler,
		},
	}
}

func (UserHandler) GetRoutes() map[string]string {
	return map[string]string{
		"user-profile": "/profile",
	}
}

func (UserHandler) GetTemplates() map[string]string {
	return map[string]string{
		"user-profile": "templates/user/profile.gohtml",
	}
}

func ProfileHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	t.TemplateOptions = framework.TemplateOptions{
		Name:  "user-profile",
		Title: "Profile",
	}
	return http.StatusOK, nil
}
