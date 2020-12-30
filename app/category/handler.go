package category

import (
	"bitbucket.org/zanvd/accountant/utility"
	"database/sql"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
)

const BaseUrl = "/category/"

func AddHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
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

		if err := InsertCategory(db, category); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/add.gohtml"))
	if err := templates.ExecuteTemplate(w, "base.gohtml", new(struct{})); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func DeleteHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	if err = DeleteCategory(db, id); err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func EditHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	category, err := GetCategory(db, id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
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

		if err = UpdateCategory(db, category); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}

	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/edit.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", category); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ListHandler(db *sql.DB, w http.ResponseWriter, _ *http.Request) (int, error) {
	categories, err := GetCategories(db)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/index.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", categories); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ViewHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	category, err := GetCategory(db, id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	templates := template.Must(template.ParseFiles("templates/base.gohtml", "templates/category/view.gohtml"))
	if err = templates.ExecuteTemplate(w, "base.gohtml", category); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
