package main

import (
	"log"
	"net/http"
	"os"

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
	if os.Getenv("ENV") == "dev" {
		sm.Env = framework.Dev
	} else {
		sm.Env = framework.Prod
	}

	r := &framework.Routes{BaseUrl: os.Getenv("BASE_URL"), Uris: make(map[string]string)}

	tb := framework.NewTemplateBuilder()
	tb.AddTemplates(map[string]string{"error": "templates/system/error.gohtml"})

	framework.RegisterHandlers(db, auth.AuthHandler{}, r, sm, tb)
	framework.RegisterRoutes(auth.AuthHandler{}, r)
	framework.RegisterTemplates(tb, auth.AuthHandler{})

	framework.RegisterHandlers(db, category.CategoryHandler{}, r, sm, tb)
	framework.RegisterRoutes(category.CategoryHandler{}, r)
	framework.RegisterTemplates(tb, category.CategoryHandler{})

	framework.RegisterHandlers(db, dashboard.DashboardHandler{}, r, sm, tb)
	framework.RegisterRoutes(dashboard.DashboardHandler{}, r)
	framework.RegisterTemplates(tb, dashboard.DashboardHandler{})

	framework.RegisterHandlers(db, public.PublicHandler{}, r, sm, tb)
	framework.RegisterRoutes(public.PublicHandler{}, r)
	framework.RegisterTemplates(tb, public.PublicHandler{})

	framework.RegisterHandlers(db, transaction.TransactionHandler{}, r, sm, tb)
	framework.RegisterRoutes(transaction.TransactionHandler{}, r)
	framework.RegisterTemplates(tb, transaction.TransactionHandler{})

	framework.RegisterHandlers(db, transaction_template.TransactionTemplateHandler{}, r, sm, tb)
	framework.RegisterRoutes(transaction_template.TransactionTemplateHandler{}, r)
	framework.RegisterTemplates(tb, transaction_template.TransactionTemplateHandler{})

	framework.RegisterHandlers(db, user.UserHandler{}, r, sm, tb)
	framework.RegisterRoutes(user.UserHandler{}, r)
	framework.RegisterTemplates(tb, user.UserHandler{})

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}
