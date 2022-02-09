package category

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/utility"
)

const BaseUrl = "/category"

type CategoryHandler struct{}

func (CategoryHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		BaseUrl: {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ListHandler,
		},
		BaseUrl + "/add": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: AddHandler,
		},
		BaseUrl + "/delete/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: DeleteHandler,
		},
		BaseUrl + "/edit/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: EditHandler,
		},
		BaseUrl + "/view/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ViewHandler,
		},
	}
}

func (CategoryHandler) GetTemplates() map[string]string {
	return map[string]string{
		"category-add":  "templates/category/add.gohtml",
		"category-edit": "templates/category/edit.gohtml",
		"category-list": "templates/category/index.gohtml",
		"category-view": "templates/category/view.gohtml",
	}
}

func AddHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		category := Category{
			UserId: t.Session.Data.User.Id,
		}
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

		if err := InsertCategory(t.DB, category); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}
	t.TemplateOptions.Name = "category-add"
	return http.StatusOK, nil
}

func DeleteHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	if err = DeleteCategory(t.DB, id, t.Session.Data.User.Id); err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}

	http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func EditHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	category, err := GetCategory(t.DB, id, t.Session.Data.User.Id)
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

		if err = UpdateCategory(t.DB, category); err != nil {
			return utility.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, BaseUrl, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}
	t.TemplateOptions = framework.TemplateOptions{
		Data: category,
		Name: "category-edit",
	}
	return http.StatusOK, nil
}

func ListHandler(t *framework.Tools, w http.ResponseWriter, _ *http.Request) (int, error) {
	categories, err := GetCategories(t.DB, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	t.TemplateOptions = framework.TemplateOptions{
		Data: categories,
		Name: "category-list",
	}
	return http.StatusOK, nil
}

func ViewHandler(t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	category, err := GetCategory(t.DB, id, t.Session.Data.User.Id)
	if err != nil {
		return utility.MapMySQLErrorToHttpCode(err), err
	}
	t.TemplateOptions = framework.TemplateOptions{
		Data: category,
		Name: "category-view",
	}
	return http.StatusOK, nil
}
