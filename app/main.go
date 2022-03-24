package main

import (
	"log"
	"net/http"

	"bitbucket.org/zanvd/accountant/auth"
	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/dashboard"
	"bitbucket.org/zanvd/accountant/framework"
	"bitbucket.org/zanvd/accountant/public"
	"bitbucket.org/zanvd/accountant/transaction"
	"bitbucket.org/zanvd/accountant/transaction_template"
	"bitbucket.org/zanvd/accountant/user"
)

func main() {
	c, err := framework.NewConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	db, err := framework.InitDatabase(c)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	cm := &framework.CacheManager{}
	cm.Connect(c)

	sm := framework.NewSessionManager(cm, c)

	m := framework.NewMailer(c)

	r := &framework.Routes{BaseUrl: c.BaseUrl, Uris: make(map[string]string)}

	tb := framework.NewTemplateBuilder()
	tb.AddTemplates(framework.GetBaseTemplates(), map[string]string{"error": "templates/system/error.gohtml"})

	mhs := []framework.ModuleHandler{
		auth.AuthHandler{},
		category.CategoryHandler{},
		dashboard.DashboardHandler{},
		public.PublicHandler{},
		transaction.TransactionHandler{},
		transaction_template.TransactionTemplateHandler{},
		user.UserHandler{},
	}
	framework.RegisterHandlers(cm, db, m, mhs, r, sm, tb)
	framework.RegisterMailTemplates([]framework.MailerConsumer{auth.AuthHandler{}}, tb)
	framework.RegisterRoutes(mhs, r)
	framework.RegisterTemplates(mhs, tb)

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}
