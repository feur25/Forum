package main

import (
	"fmt"
	"log"

	"net/http"
)

func handleGetCode(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "code"
	data.Page.Style = "code"
	code := r.FormValue("code")
	if r.FormValue("code_drop") == "getcode" {
		send_reset_mail()
		log.Print("\ncode : " + data.code + "\n")
	}
	if code == data.code && len(data.code) >= 5 {
		data.Authorized = true
		if data.Authorized == true {
			http.Redirect(w, r, "http://"+Host+":"+Port+"/update", http.StatusMovedPermanently)
		}
	}
	tmpl.ExecuteTemplate(w, "code", data)
}
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if data.Authorized == true {
		data.Page.Title = "Update"
		data.Page.Style = "update"

		password := r.FormValue("password_remove")
		passwordcheck := r.FormValue("checkpassword")
		if password != passwordcheck {
			var bad_password = `<p style="color: red;">Les deux mots de passe ne correspondent pas `
			fmt.Fprintf(w, bad_password)
		}
		log.Print("\n" + passwordcheck)
		log.Print(password)
		if r.FormValue("envoyer") == "Envoyer" {
			if password == passwordcheck {
				log.Print("Le mot de passe a bien changer")
				updateUserInDB(w, r, MD5(password))
				http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
			}
			if password == data.Auth.password {
				var password_active = `<p style="color: red;">Le nouveau mot de passe resemble a votre ancien, mot de passe `
				fmt.Fprintf(w, password_active)
			}
		}
	} else {
		var bad_code_message = `<p style="color: red;">Vous devez d'abord récuperer votre code de vérification, avant de vouloir rinitialiser votre mot de passe !`
		fmt.Fprintf(w, bad_code_message)
	}
	tmpl.ExecuteTemplate(w, "update", data)
}
