package main

import (
	"log"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Home"
	data.Page.Style = "home"

	data.Page.Topics = getMostRecentTopics(20)

	log.Print(data.Page.Topics)
	tmpl.ExecuteTemplate(w, "home", data)
}
