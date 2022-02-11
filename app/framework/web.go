package framework

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/zanvd/accountant/utility"
	"github.com/go-sql-driver/mysql"
)

type Middleware interface {
	PreRequest(t *Tools, r *http.Request) error
	PreRender(t *Tools, w http.ResponseWriter) error
	Render(t *Tools, w http.ResponseWriter) error
	PostRequest(t *Tools, w http.ResponseWriter) error
}

type Routes struct {
	BaseUrl string
	Uris    map[string]string
}

type Tools struct {
	DB              *sql.DB
	Routes          *Routes
	Session         Session
	TemplateOptions TemplateOptions
}

type Endpoint struct {
	Auth    AuthSettings
	Handler func(t *Tools, w http.ResponseWriter, r *http.Request) (int, error)
}

type ModuleHandler interface {
	GetHandlers() map[string]Endpoint // path = endpoint handler
	GetRoutes() map[string]string
	GetTemplates() map[string]string
}

// TODO Rework session so it is managed through SessionManager from handlers.
type AppHandler struct {
	Endpoint        Endpoint
	SessionManager  *SessionManager
	TemplateBuilder *TemplateBuilder
	Tools           *Tools
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Have to obtain a fresh set of Tools on each request. It would otherwise persist for the same handler
	// resulting in things like error message and status being left over in the TemplateOptions.
	ah.Tools = newTools(ah.Tools.DB, ah.Tools.Routes)
	// Obtain session.
	if err := ah.SessionManager.GetSession(ah.Tools, r); err != nil {
		// TODO Log error.
	}
	// Check authorization.
	if !ah.Endpoint.Auth.Public && !ah.Tools.Session.Data.LoggedIn {
		ah.handleErrors(errors.New("unauthorized"), http.StatusUnauthorized, w, r)
		return
	}
	// Handle request.
	status, err := ah.Endpoint.Handler(ah.Tools, w, r)
	if err != nil {
		ah.handleErrors(err, status, w, r)
		return
	}
	// Clear or update session.
	if ah.Tools.Session.Clear {
		if err := ah.SessionManager.ClearSession(ah.Tools, w); err != nil {
			ah.handleErrors(err, http.StatusInternalServerError, w, r)
		}
		return
	}
	if err := ah.SessionManager.WriteSession(ah.Tools, w); err != nil {
		ah.handleErrors(err, http.StatusInternalServerError, w, r)
		return
	}
	if status == http.StatusTemporaryRedirect {
		return
	}
	// Render template.
	if err := ah.TemplateBuilder.Render(ah.Tools, w); err != nil {
		ah.handleErrors(err, http.StatusInternalServerError, w, r)
		return
	}
}

func (ah AppHandler) handleErrors(err error, status int, w http.ResponseWriter, r *http.Request) {
	log.Printf("error occurred (%s): %+v", r.URL.RequestURI(), err)

	// Handle DB errors with a custom text.
	message := err.Error()
	if _, ok := err.(*mysql.MySQLError); ok {
		message = utility.GetMySQLErrorMessage(err)
	} else if err == sql.ErrConnDone || err == sql.ErrTxDone {
		message = "Something went wrong."
		status = http.StatusInternalServerError
	} else if err == sql.ErrNoRows {
		message = "The requested data hasn't been found."
		status = http.StatusNotFound
	}

	w.WriteHeader(status)

	ah.Tools.TemplateOptions.ErrorMessage = message
	ah.Tools.TemplateOptions.ErrorStatus = status
	if ah.Tools.TemplateOptions.Name == "" {
		ah.Tools.TemplateOptions.Name = "error"
	}
	if ah.Tools.TemplateOptions.Title == "" {
		ah.Tools.TemplateOptions.Title = fmt.Sprintf("Error (%d)", status)
	}
	if err := ah.TemplateBuilder.Render(ah.Tools, w); err != nil {
		// TODO Handle error.
		return
	}
	/*templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/system/error.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/
}

func RegisterHandlers(db *sql.DB, mh ModuleHandler, r *Routes, sm *SessionManager, tb *TemplateBuilder) {
	ehs := mh.GetHandlers()
	for p, eh := range ehs {
		http.Handle(p, newAppHandler(eh, sm, tb, newTools(db, r)))
	}
}

func RegisterRoutes(mh ModuleHandler, r *Routes) {
	for n, u := range mh.GetRoutes() {
		if _, ok := r.Uris[n]; ok {
			log.Panicln("error: URI already added with name", n)
		}
		r.Uris[n] = u
	}
}

func RegisterTemplates(tb *TemplateBuilder, mh ModuleHandler) {
	tb.AddTemplates(mh.GetTemplates())
}

func newAppHandler(e Endpoint, sm *SessionManager, tb *TemplateBuilder, t *Tools) AppHandler {
	return AppHandler{
		Endpoint:        e,
		SessionManager:  sm,
		TemplateBuilder: tb,
		Tools:           t,
	}
}

func newTools(db *sql.DB, r *Routes) *Tools {
	return &Tools{
		DB:              db,
		Routes:          r,
		Session:         Session{},
		TemplateOptions: TemplateOptions{},
	}
}
