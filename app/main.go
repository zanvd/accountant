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
	"bitbucket.org/zanvd/accountant/utility"
)

func main() {
	db, err := utility.InitDatabase()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	var sm *framework.SessionManager = new(framework.SessionManager)
	sm.Connect(framework.ConnectionConfig{
		Db:       0,
		Host:     "redis",
		Password: "",
		Port:     "6379",
		Username: "",
	})

	tb := framework.NewTemplateBuilder()
	tb.AddTemplates(map[string]string{"error": "templates/system/error.gohtml"})

	framework.RegisterHandlers(db, auth.AuthHandler{}, sm, tb)
	framework.RegisterTemplates(tb, auth.AuthHandler{})

	framework.RegisterHandlers(db, category.CategoryHandler{}, sm, tb)
	framework.RegisterTemplates(tb, category.CategoryHandler{})

	framework.RegisterHandlers(db, dashboard.DashboardHandler{}, sm, tb)
	framework.RegisterTemplates(tb, dashboard.DashboardHandler{})

	framework.RegisterHandlers(db, public.PublicHandler{}, sm, tb)
	framework.RegisterTemplates(tb, public.PublicHandler{})

	framework.RegisterHandlers(db, transaction.TransactionHandler{}, sm, tb)
	framework.RegisterTemplates(tb, transaction.TransactionHandler{})

	framework.RegisterHandlers(db, transaction_template.TransactionTemplateHandler{}, sm, tb)
	framework.RegisterTemplates(tb, transaction_template.TransactionTemplateHandler{})

	framework.RegisterHandlers(db, user.UserHandler{}, sm, tb)
	framework.RegisterTemplates(tb, user.UserHandler{})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}
