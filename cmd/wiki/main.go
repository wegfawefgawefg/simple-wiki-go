package main

import (
	"html/template"
	"log"
	"net/http"
)

const saved_pages_path = "pages"
const templates_path = "templates"

type Page struct {
	Title         string
	InternalTitle string
	Body          []byte
}

var funcMap = template.FuncMap{"nl2br": nl2br}

var templates = template.Must(
	template.New("").Funcs(funcMap).ParseFiles(
		templates_path+"/edit.html",
		templates_path+"/view.html",
		templates_path+"/index.html",
		templates_path+"/add_new_page.html",
	))

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/add_new_page", addNewPageHandler)
	http.HandleFunc("/internal_add_new_page", internalAddNewPageHandler)
	http.HandleFunc("/index", indexHandler)
	// http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
