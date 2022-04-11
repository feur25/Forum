package main

import (
	"fmt"
	"log"

	"net/http"
)

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	send_reset_mail()

	data.Page.Title = "Update"
	data.Page.Style = "update"

	password := r.FormValue("password_remove")
	passwordcheck := r.FormValue("checkpassword")
	codecheck := r.FormValue("code")
	log.Print("\ntest\n" + data.code + codecheck)
	if password != passwordcheck {
		var bad_password = `<p style="color: red;">Les deux mots de passe ne correspondent pas <\p>`
		fmt.Fprintf(w, bad_password)
	}
	log.Print("\n" + passwordcheck)
	log.Print(password)
	if r.FormValue("envoyer") == "Envoyer" {
		if password == passwordcheck {
			log.Print("Le mot de passe a bien changer")
			updateUserInDB(w, r, MD5(password))
			http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
		} else {
			log.Print("Il y a une erreur ")
		}
		if password == data.Auth.password {
			var password_active = `<p style="color: red;">Le nouveau mot de passe resemble a votre ancien, mot de passe <\p>`
			fmt.Fprintf(w, password_active)
		}
	}
	tmpl.ExecuteTemplate(w, "update", data)
}
