package category

import (
	"database/sql"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
)

const BaseUrl = "/category/"

type AddHandler struct {
	Database *sql.DB
}

func (ah AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		category := Category{}
		if color := strings.TrimSpace(r.FormValue("color")); color != "" {
			category.Color = color
		} else {
			category.Color = defaultColor
		}
		category.Description = strings.TrimSpace(r.FormValue("description"))
		if name := strings.TrimSpace(r.FormValue("name")); name != "" {
			category.Name = name
		}
		if textColor := strings.TrimSpace(r.FormValue("text-color")); textColor != "" {
			category.TextColor = textColor
		} else {
			category.TextColor = defaultTextColor
		}

		if err := InsertCategory(ah.Database, category); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	}

	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/add.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", new(struct{})); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type DeleteHandler struct {
	Database *sql.DB
}

func (dh DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err = DeleteCategory(dh.Database, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
}

type EditHandler struct {
	Database *sql.DB
}

func (eh EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	category, err := GetCategory(eh.Database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if r.Method == "POST" {
		if color := strings.TrimSpace(r.FormValue("color")); color != "" {
			category.Color = color
		} else {
			category.Color = defaultColor
		}
		category.Description = strings.TrimSpace(r.FormValue("description"))
		if name := strings.TrimSpace(r.FormValue("name")); name != "" {
			category.Name = name
		}
		if textColor := strings.TrimSpace(r.FormValue("text-color")); textColor != "" {
			category.TextColor = textColor
		} else {
			category.TextColor = defaultTextColor
		}

		if err := UpdateCategory(eh.Database, category); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	}

	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/edit.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ListHandler struct {
	Database *sql.DB
}

func (lh ListHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	categories, err := GetCategories(lh.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/index.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ViewHandler struct {
	Database *sql.DB
}

func (vh ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	category, err := GetCategory(vh.Database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/view.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
