package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/user"
)

type AuthHandler struct{}

func (AuthHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		"/confirm-account": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: ConfirmAccountHandler,
		},
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
		"/new-confirm-account": {
			Auth: framework.AuthSettings{
				Public: true,
			},
			Handler: NewConfirmationHandler,
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

func (AuthHandler) GetMailTemplates() map[string]string {
	return map[string]string{
		"auth-mail-new-confirmation": "templates/mail/auth/new-confirmation.gohtml",
		"auth-mail-reset-password":   "templates/mail/auth/reset-password.gohtml",
		"auth-mail-welcome":          "templates/mail/auth/welcome.gohtml",
	}
}

func (AuthHandler) GetRoutes() map[string]string {
	return map[string]string{
		"auth-confirm-account": "/confirm-account",
		"auth-forgot-password": "/forgot-password",
		"auth-login":           "/login",
		"auth-logout":          "/logout",
		"auth-new-confirm":     "/new-confirm-account",
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

func ConfirmAccountHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	username, token := "", ""
	if username = strings.TrimSpace(r.URL.Query().Get("username")); username == "" {
		return http.StatusBadRequest, errors.New("missing username. Please request a new confirmation link")
	}
	if token = strings.TrimSpace(r.URL.Query().Get("token")); token == "" {
		return http.StatusBadRequest, errors.New("missing token. Please request a new confirmation link")
	}
	cu, err := readToken(t.CacheManager, "auth-confirm", token)
	if err != nil {
		return http.StatusInternalServerError, err
	} else if cu == "" || username != cu {
		return http.StatusBadRequest, errors.New("invalid token. Please request a new confirmation link")
	}
	u, err := user.GetUserByUsername(t.DB, username)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid username. Please request a new confirmation link")
	}
	u.Confirmed = true
	if err := user.UpdateUser(t.DB, u); err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}
	deleteToken(t.CacheManager, "auth-confirm", token)
	http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["auth-login"], http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func ForgotPasswordHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	sent := false
	if r.Method == "POST" {
		un := ""
		if un = strings.TrimSpace(r.PostFormValue("username")); un == "" {
			return http.StatusBadRequest, errors.New("please provide a username")
		}
		u, err := user.GetUserByUsername(t.DB, un)
		if err != nil {
			return http.StatusBadRequest, errors.New("invalid username")
		}
		token := createToken()
		if err := cacheToken(t.CacheManager, un, "auth-pass-reset", token, time.Hour*24); err != nil {
			return http.StatusInternalServerError, err
		}
		if err := sendEmail(*rd, "Reset password", "auth-mail-reset-password", token, t, u); err != nil {
			return http.StatusInternalServerError, err
		}
		sent = true
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Sent bool
		}{
			Sent: sent,
		},
		Name:  "auth-forgot-password",
		Title: "Forgot Password",
	}
	return http.StatusOK, nil
}

func LoginHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
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
		if !u.Confirmed {
			return http.StatusForbidden, errors.New("please confirm your email")
		}
		rd.Session.Data.LoggedIn = true
		rd.Session.Data.User = framework.SessionUser{
			Id:       u.Id,
			Email:    u.Email,
			Username: u.Username,
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["dashboard"], http.StatusSeeOther)
		return http.StatusSeeOther, nil
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Name:  "auth-login",
		Title: "Login",
	}
	return http.StatusOK, nil
}

func LogoutHanlder(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if err := t.SessionManager.ClearSession(&rd.Session, w); err != nil {
		return http.StatusInternalServerError, nil
	}
	http.Redirect(w, r, t.Routes.BaseUrl, http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func NewConfirmationHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	sent := false
	if r.Method == "POST" {
		un := ""
		if un = strings.TrimSpace(r.PostFormValue("username")); un == "" {
			return http.StatusBadRequest, errors.New("please provide a username")
		}
		u, err := user.GetUserByUsername(t.DB, un)
		if err != nil {
			return http.StatusBadRequest, errors.New("invalid username")
		}
		if u.Confirmed {
			return http.StatusBadRequest, errors.New("you have already confirmed your account")
		}
		token := createToken()
		if err := cacheToken(t.CacheManager, un, "auth-confirm", token, time.Hour*24); err != nil {
			return http.StatusInternalServerError, err
		}
		if err := sendEmail(*rd, "Confirm account", "auth-mail-new-confirmation", token, t, u); err != nil {
			return http.StatusInternalServerError, err
		}
		sent = true
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Sent bool
		}{
			Sent: sent,
		},
		Name:  "auth-new-confirm",
		Title: "Resend Confirmation",
	}
	return http.StatusOK, nil
}

func PasswordResetHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	username, token := "", ""
	if username = strings.TrimSpace(r.URL.Query().Get("username")); username == "" {
		return http.StatusBadRequest, errors.New("missing username. Please request a new password reset link")
	}
	if token = strings.TrimSpace(r.URL.Query().Get("token")); token == "" {
		return http.StatusBadRequest, errors.New("missing token. Please request a new password reset link")
	}
	cu, err := readToken(t.CacheManager, "auth-pass-reset", token)
	if err != nil {
		return http.StatusInternalServerError, err
	} else if cu == "" || username != cu {
		return http.StatusBadRequest, errors.New("invalid token. Please request a new password reset link")
	}
	u, err := user.GetUserByUsername(t.DB, username)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid username. Please request a new password reset link")
	}

	if r.Method == "POST" {
		p, pr := "", ""
		if p = strings.TrimSpace(r.PostFormValue("password")); p == "" {
			return http.StatusBadRequest, errors.New("please provide a password")
		}
		if pr = strings.TrimSpace(r.PostFormValue("password-repeat")); pr == "" {
			return http.StatusBadRequest, errors.New("please provide a repeated password")
		}
		ph, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		u.Password = string(ph)
		if err := user.UpdateUser(t.DB, u); err != nil {
			return http.StatusInternalServerError, err
		}
		deleteToken(t.CacheManager, "auth-pass-reset", token)
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["auth-login"], http.StatusSeeOther)
		return http.StatusSeeOther, nil
	}

	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Token string
			User  user.User
		}{
			Token: token,
			User:  u,
		},
		Name:  "auth-password-reset",
		Title: "Password Reset",
	}
	return http.StatusOK, nil
}

func RegisterHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
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
			return framework.MapMySQLErrorToHttpCode(err), err
		}
		token := createToken()
		if err := cacheToken(t.CacheManager, un, "auth-confirm", token, time.Hour*24); err != nil {
			return http.StatusInternalServerError, err
		}
		if err := sendEmail(*rd, "Welcome to Accountant", "auth-mail-welcome", token, t, u); err != nil {
			return http.StatusInternalServerError, err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["auth-login"], http.StatusSeeOther)
		return http.StatusSeeOther, nil
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Name:  "auth-register",
		Title: "Register",
	}
	return http.StatusOK, nil
}

func cacheToken(cm *framework.CacheManager, username string, prefix string, token string, ttl time.Duration) error {
	return cm.Set(prefix+":"+token, username, ttl)
}

func createToken() string {
	return uuid.New().String()
}

func deleteToken(cm *framework.CacheManager, prefix string, token string) (int64, error) {
	return cm.Delete(prefix + ":" + token)
}

func readToken(cm *framework.CacheManager, prefix string, token string) (string, error) {
	return cm.Get(prefix + ":" + token)
}

func sendEmail(rd framework.RequestData, subject string, tmplName string, token string, tools *framework.Tools, u user.User) (err error) {
	rd.TemplateOptions = framework.TemplateOptions{
		Data: struct {
			Token string
			User  user.User
		}{
			Token: token,
			User:  u,
		},
		Name: tmplName,
	}
	m := framework.Mail{
		From:    tools.Mailer.DefaultFrom,
		Subject: subject,
		To:      []string{u.Email},
	}
	if err = m.RenderBody(tools.Routes, &rd, tools.TemplateBuilder); err != nil {
		return err
	}
	if err = tools.Mailer.Send(m); err != nil {
		return err
	}
	return
}
