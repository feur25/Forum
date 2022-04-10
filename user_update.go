package main

import (
	"fmt"
	"log"

	"net/http"
)

func handleUpdate(w http.ResponseWriter, r *http.Request) {

	data.Page.Title = "Update"
	data.Page.Style = "update"

	password := r.FormValue("password_remove")
	passwordcheck := r.FormValue("checkpassword")
	if password != passwordcheck {
		var bad_password = `<p style="color: red;">Les deux mots de passe ne correspondent pas <\p>`
		fmt.Fprintf(w, bad_password)
	}
	log.Print("\n" + passwordcheck)
	log.Print(password)
	if r.FormValue("submit") == "Envoyer" {
		if password == passwordcheck {
			updateUserInDB(w, r, password)
		}
	}
	tmpl.ExecuteTemplate(w, "update", data)
}
