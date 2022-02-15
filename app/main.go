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
	framework.RegisterHandlers(db, mhs, r, sm, tb)
	framework.RegisterRoutes(mhs, r)
	framework.RegisterTemplates(mhs, tb)

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}
