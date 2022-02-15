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

type RequestData struct {
	Session         Session
	TemplateOptions TemplateOptions
}

type Routes struct {
	BaseUrl string
	Uris    map[string]string
}

type Tools struct {
	DB              *sql.DB
	Routes          *Routes
	SessionManager  *SessionManager
	TemplateBuilder *TemplateBuilder
}

type Endpoint struct {
	Auth    AuthSettings
	Handler func(rd *RequestData, t *Tools, w http.ResponseWriter, r *http.Request) (int, error)
}

type ModuleHandler interface {
	GetHandlers() map[string]Endpoint // path = endpoint handler
	GetRoutes() map[string]string     // name = route
	GetTemplates() map[string]string  // name = path
}

type AppHandler struct {
	Endpoint    Endpoint
	RequestData RequestData
	Tools       *Tools
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	ah.RequestData = RequestData{}
	// Obtain session.
	if ah.RequestData.Session, err = ah.Tools.SessionManager.GetSession(r); err != nil {
		// TODO: Log error.
	}
	// Check authorization.
	if !ah.Endpoint.Auth.Public && !ah.RequestData.Session.Data.LoggedIn {
		ah.handleErrors(errors.New("unauthorized"), http.StatusUnauthorized, w, r)
		return
	}
	// Handle request.
	status, err := ah.Endpoint.Handler(&ah.RequestData, ah.Tools, w, r)
	if err != nil {
		ah.handleErrors(err, status, w, r)
		return
	} else if status == http.StatusTemporaryRedirect {
		return
	}

	// Update session.
	if err := ah.Tools.SessionManager.WriteSession(ah.RequestData.Session, w); err != nil {
		ah.handleErrors(err, http.StatusInternalServerError, w, r)
		return
	}
	// Render template.
	if err := ah.Tools.TemplateBuilder.Render(ah.Tools.Routes, &ah.RequestData, ah.Tools, w); err != nil {
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

	ah.RequestData.TemplateOptions.ErrorMessage = message
	ah.RequestData.TemplateOptions.ErrorStatus = status
	if ah.RequestData.TemplateOptions.Name == "" {
		ah.RequestData.TemplateOptions.Name = "error"
	}
	if ah.RequestData.TemplateOptions.Title == "" {
		ah.RequestData.TemplateOptions.Title = fmt.Sprintf("Error (%d)", status)
	}
	if err := ah.Tools.TemplateBuilder.Render(ah.Tools.Routes, &ah.RequestData, ah.Tools, w); err != nil {
		// TODO: Handle error.
		return
	}
	/*templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/system/error.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/
}

func RegisterHandlers(db *sql.DB, mhs []ModuleHandler, r *Routes, sm *SessionManager, tb *TemplateBuilder) {
	for _, mh := range mhs {
	ehs := mh.GetHandlers()
		t := &Tools{
			DB:              db,
			Routes:          r,
			SessionManager:  sm,
			TemplateBuilder: tb,
		}
	for p, eh := range ehs {
			http.Handle(p, AppHandler{Endpoint: eh, Tools: t})
		}
	}
}

func RegisterRoutes(mhs []ModuleHandler, r *Routes) {
	for _, mh := range mhs {
	for n, u := range mh.GetRoutes() {
		if _, ok := r.Uris[n]; ok {
			log.Panicln("error: URI already added with name", n)
		}
		r.Uris[n] = u
	}
}
}

func RegisterTemplates(mhs []ModuleHandler, tb *TemplateBuilder) {
	for _, mh := range mhs {
		tb.AddTemplates(baseTmpls, mh.GetTemplates())
	}
}
