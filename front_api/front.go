package front_api

import (
	"html/template"
	"net/http"
)

func handleDefault(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./htmx/index.html"))
	tmpl.Execute(w, nil)
}

func HandleEndPoints() {
	http.HandleFunc("/", handleDefault)
}
