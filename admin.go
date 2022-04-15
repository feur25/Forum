package main

import (
	"fmt"
	"log"
	"net/http"
)

func new_tag(name string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `tags` (`name`) VALUES ('%s')", name))
	log.Print("Le tag : [" + name + "] a bien été implémenté")
	checkError(err)
	defer insert.Close()
	return nil
}
func handleAdminPanel(w http.ResponseWriter, r *http.Request) {
	log.Print(data.Auth.Admin)
	data.Page.Title = "Admin"
	data.Page.Style = "admin"
	tag := r.FormValue("tag")
	if data.Auth.Admin == 1 {
		if r.FormValue("envoyer") == "Envoyer" {
			log.Print("let's go !")
			new_tag(tag)
		}
	}
	if data.Auth.Admin == 0 {
		go log.Print("L'utilisateur n'es pas un admin")
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusFound)
	}
	tmpl.ExecuteTemplate(w, "admin", data)
}
