package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func mainhandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/main.html")
	t.Execute(w, nil)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainhandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":1234", r)

}
