package category

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"bitbucket.org/zanvd/accountant/framework"
)

const BaseUri = "/category"

type CategoryHandler struct{}

func (CategoryHandler) GetHandlers() map[string]framework.Endpoint {
	return map[string]framework.Endpoint{
		BaseUri: {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ListHandler,
		},
		BaseUri + "/add": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: AddHandler,
		},
		BaseUri + "/delete/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: DeleteHandler,
		},
		BaseUri + "/edit/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: EditHandler,
		},
		BaseUri + "/view/": {
			Auth: framework.AuthSettings{
				Public: false,
			},
			Handler: ViewHandler,
		},
	}
}

func (CategoryHandler) GetRoutes() map[string]string {
	return map[string]string{
		"category-add":    BaseUri + "/add",
		"category-delete": BaseUri + "/delete",
		"category-edit":   BaseUri + "/edit",
		"category-list":   BaseUri,
		"category-view":   BaseUri + "/view",
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

func AddHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		category := Category{
			UserId: rd.Session.Data.User.Id,
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
			return framework.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["category-list"], http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Name:  "category-add",
		Title: "Add Category",
	}
	return http.StatusOK, nil
}

func DeleteHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	if err = DeleteCategory(t.DB, id, rd.Session.Data.User.Id); err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}

	http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["category-list"], http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}

func EditHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	category, err := GetCategory(t.DB, id, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
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
			return framework.MapMySQLErrorToHttpCode(err), err
		}
		http.Redirect(w, r, t.Routes.BaseUrl+t.Routes.Uris["category-list"], http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect, nil
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data:  category,
		Name:  "category-edit",
		Title: "Edit Category",
	}
	return http.StatusOK, nil
}

func ListHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, _ *http.Request) (int, error) {
	categories, err := GetCategories(t.DB, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data:  categories,
		Name:  "category-list",
		Title: "Categories",
	}
	return http.StatusOK, nil
}

func ViewHandler(rd *framework.RequestData, t *framework.Tools, w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, err
	}
	category, err := GetCategory(t.DB, id, rd.Session.Data.User.Id)
	if err != nil {
		return framework.MapMySQLErrorToHttpCode(err), err
	}
	rd.TemplateOptions = framework.TemplateOptions{
		Data:  category,
		Name:  "category-view",
		Title: "View Category",
	}
	return http.StatusOK, nil
}
