package framework

import (
	"log"
	"net/http"
	"text/template"

	"bitbucket.org/zanvd/accountant/convert"
)

var sysTmpls = [3]string{
	"templates/base.gohtml",
	"templates/system/error.gohtml",
	"templates/system/form_errors.gohtml",
}

type TemplateBuilder struct {
	funcs    template.FuncMap
	tmplList map[string]*template.Template // [name] = Template
}

type TemplateData struct {
	Data         interface{}
	ErrorMessage string
	ErrorStatus  int
	Routes       *Routes
	Session      SessionData
	Title        string
}

type TemplateOptions struct {
	Name         string
	Data         interface{}
	ErrorMessage string
	ErrorStatus  int
	Title        string
}

func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{
		funcs: template.FuncMap{
			"dbToDisplayDate": func(dbDate string) string {
				return convert.DbToDisplayDate(dbDate)
			},
			"today": func() string {
				return convert.CurrentDateInDisplayFormat()
			},
			"url": func(name string, r *Routes) string {
				uri, ok := r.Uris[name]
				if !ok {
					log.Panicln("error: no route with name", name)
				}
				return r.BaseUrl + uri
			},
		},
		tmplList: make(map[string]*template.Template),
	}
}

func (tb *TemplateBuilder) AddTemplates(tmpls map[string]string) {
	for n, p := range tmpls {
		if _, ok := tb.tmplList[n]; ok {
			log.Panicln("error: template already added with name", n)
		}
		tb.tmplList[n] = tb.load(n, p)
	}
}

func (tb *TemplateBuilder) Render(t *Tools, w http.ResponseWriter) error {
	if t.TemplateOptions.Data == nil {
		t.TemplateOptions.Data = new(struct{})
	}
	d := TemplateData{
		Data:         t.TemplateOptions.Data,
		Title:        t.TemplateOptions.Title,
		ErrorMessage: t.TemplateOptions.ErrorMessage,
		ErrorStatus:  t.TemplateOptions.ErrorStatus,
		Routes:       t.Routes,
		Session:      t.Session.Data,
	}
	return tb.tmplList[t.TemplateOptions.Name].ExecuteTemplate(w, "base.gohtml", d)
}

func (tb *TemplateBuilder) load(n string, p string) *template.Template {
	paths := []string{p}
	for _, sp := range sysTmpls {
		paths = append(paths, sp)
	}
	return template.Must(template.New(n).Funcs(tb.funcs).ParseFiles(paths...))
}
