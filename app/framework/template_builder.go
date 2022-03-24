package framework

import (
	"io"
	"log"
	"text/template"

	"bitbucket.org/zanvd/accountant/convert"
)

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
	Data         interface{}
	ErrorMessage string
	ErrorStatus  int
	Name         string
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

func (tb *TemplateBuilder) AddTemplates(baseTmpls []string, tmpls map[string]string) {
	for n, p := range tmpls {
		if _, ok := tb.tmplList[n]; ok {
			log.Panicln("error: template already added with name", n)
		}
		tb.tmplList[n] = tb.load(baseTmpls, n, p)
	}
}

func (tb *TemplateBuilder) Render(r *Routes, rd *RequestData, w io.Writer) error {
	if rd.TemplateOptions.Data == nil {
		rd.TemplateOptions.Data = new(struct{})
	}
	d := TemplateData{
		Data:         rd.TemplateOptions.Data,
		Title:        rd.TemplateOptions.Title,
		ErrorMessage: rd.TemplateOptions.ErrorMessage,
		ErrorStatus:  rd.TemplateOptions.ErrorStatus,
		Routes:       r,
		Session:      rd.Session.Data,
	}
	return tb.tmplList[rd.TemplateOptions.Name].ExecuteTemplate(w, "base.gohtml", d)
}

func (tb *TemplateBuilder) load(baseTmpls []string, n string, p string) *template.Template {
	baseTmpls = append(baseTmpls, p)
	return template.Must(template.New(n).Funcs(tb.funcs).ParseFiles(baseTmpls...))
}
