package main

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/zanvd/accountant/category"
	"bitbucket.org/zanvd/accountant/transaction"
	"bitbucket.org/zanvd/accountant/transaction_template"
	"bitbucket.org/zanvd/accountant/user"
	"bitbucket.org/zanvd/accountant/utility"
)

func main() {
	args := os.Args[1:]
	switch args[0] {
	case "createTables":
		createTables()
	case "help":
		help()
	default:
		help()
	}
}

func createTables() {
	db, err := utility.InitDatabase()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	user.CreateUserTable(db)
	category.CreateCategoryTable(db)
	transaction.CreateTransactionsTable(db)
	transaction_template.CreateTransactionTemplateTable(db)

	fmt.Println("Done.")
}

func help() {
	fmt.Println(`Available commands:
	createTables	Creates DB tables if they don't exist.
	help			Prints this help text.`)
}
