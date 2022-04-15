package main

import (
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Home"
	data.Page.Style = "home"
	db.QueryRow("topic").Scan(&data.Topic)
	tmpl.ExecuteTemplate(w, "home", data)
}
