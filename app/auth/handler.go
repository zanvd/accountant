package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/user"
	"bitbucket.org/zanvd/accountant/utility"
)

type AuthHandler struct{}

func (AuthHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		"/forgot-password": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: ForgotPasswordHandler,
		},
		"/login": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: LoginHandler,
		},
		"/logout": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: LogoutHanlder,
		},
		"/password-reset": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: PasswordResetHandler,
		},
		"/register": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: RegisterHandler,
		},
	}
}

func (AuthHandler) GetRoutes() map[string]string {
	return map[string]string{
		"auth-forgot-password": "/forgot-password",
		"auth-login":           "/login",
		"auth-logout":          "/logout",
		"auth-password-reset":  "/password-reset",
		"auth-register":        "/register",
	}
}

func (AuthHandler) GetTemplates() map[string]string {
	return map[string]string{
		"auth-forgot-password": "templates/auth/forgot-password.gohtml",
		"auth-login":           "templates/auth/login.gohtml",
		"auth-password-reset":  "templates/auth/password-reset.gohtml",
		"auth-register":        "templates/auth/register.gohtml",
	}
}

func ForgotPasswordHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {

	}
	t.TemplateOptions = framework.TemplateOptions{
		Name:  "auth-forgot-password",
		Title: "Forgot Password",
	}
	return http.StatusOK, nil
}

func LoginHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		p, un := "", ""
		if p = strings.TrimSpace(r.PostFormValue("password")); p == "" {
			return http.StatusBadRequest, errors.New("please provide a password")
		}
		if un = strings.TrimSpace(r.PostFormValue("username")); un == "" {
			return http.StatusBadRequest, errors.New("please provide a username")
		}
		u, err := user.GetUserByUsername(t.DB, un)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p)) != nil {
			return http.StatusBadRequest, errors.New("invalid username or password")
		}
		t.Session.Data.LoggedIn = true
		t.Session.Data.User = framework.SessionUser{
			Id:       u.Id,
			Email:    u.Email,
			Username: u.Username,
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["dashboard"], http.StatusTemporaryRedirect)
	}
	t.TemplateOptions = framework.TemplateOptions{
		Name:  "auth-login",
		Title: "Login",
	}
	return http.StatusOK, nil
}

func LogoutHanlder(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if err := t.SessionManager.ClearSession(&t.Session, w); err != nil {
		return http.StatusInternalServerError, nil
	}
	http.Redirect(w, r, t.Routes.BaseUrl, http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func PasswordResetHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {

	}
	t.TemplateOptions = framework.TemplateOptions{
		Name:  "auth-password-reset",
		Title: "Password Reset",
	}
	return http.StatusOK, nil
}

func RegisterHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		e, p, pr, un := "", "", "", ""
		if e = strings.TrimSpace(r.PostFormValue("email")); e == "" {
			return http.StatusBadRequest, errors.New("please provide an email")
		}
		if p = strings.TrimSpace(r.PostFormValue("password")); p == "" {
			return http.StatusBadRequest, errors.New("please provide a password")
		}
		if pr = strings.TrimSpace(r.PostFormValue("password-repeat")); pr == "" {
			return http.StatusBadRequest, errors.New("please provide a repeated password")
		}
		if p != pr {
			return http.StatusBadRequest, errors.New("provided password and repeated password do not match. Please try again")
		}
		if un = strings.TrimSpace(r.PostFormValue("username")); un == "" {
			return http.StatusBadRequest, errors.New("please provide a username")
		}
		ph, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		p = string(ph)
		u := user.User{
			Email:    e,
			Password: p,
			Username: un,
		}
		if err := user.InsertUser(t.DB, u); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["auth-login"], http.StatusTemporaryRedirect)
	}
	t.TemplateOptions = framework.TemplateOptions{
		Name:  "auth-register",
		Title: "Register",
	}
	return http.StatusOK, nil
}
