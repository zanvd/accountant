package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

const timeFormat = "15:04 02. 01. 2006"

type transaction struct {
	Id          int
	Category    string
	Date        string
	Description string
	Name        string
	Value       float64
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func editHandler(w http.ResponseWriter, _ *http.Request) {
	transaction := transaction{
		Id:          1,
		Category:    "Category 1",
		Date:        time.Now().Format(timeFormat),
		Description: "Description for transaction 1.",
		Name:        "Transaction 1",
		Value:       1234,
	}
	templates := template.Must(template.ParseFiles("templates/base.html", "templates/edit.html"))
	if err := templates.ExecuteTemplate(w, "base.html", transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	transactions := [10]transaction{}
	for i := 0; i < 10; i++ {
		transactions[i] = transaction{
			Id:          i + 1,
			Category:    "Category " + strconv.FormatInt(int64(i+1), 10),
			Date:        time.Now().Format(timeFormat),
			Description: "Some transaction " + strconv.FormatInt(int64(i+1), 10),
			Name:        "Transaction " + strconv.FormatInt(int64(i+1), 10),
			Value:       (float64)(i+1) * 3,
		}
	}
	templates := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	if err := templates.ExecuteTemplate(w, "base.html", transactions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, _ *http.Request) {
	transaction := transaction{
		Id:          1,
		Category:    "Category 1",
		Date:        time.Now().Format(timeFormat),
		Description: "Description for transaction 1.",
		Name:        "Transaction 1",
		Value:       1234,
	}
	templates := template.Must(template.ParseFiles("templates/base.html", "templates/view.html"))
	if err := templates.ExecuteTemplate(w, "base.html", transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
