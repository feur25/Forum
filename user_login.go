package main

import (
	"fmt"
	"log"

	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	data.Page.Title = "Login"
	data.Page.Style = "login"

	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Print("pseudo : " + username + " mdp : " + password)
	if IsButtonPressed(r, "login") {
		checkUserLogin(w, r, username, password)
	}
	tmpl.ExecuteTemplate(w, "login-page", data)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	defer Redirect(w, r, "/home")

	data.User = User{}
	data.Login = false

}

func checkUserLogin(w http.ResponseWriter, r *http.Request, username, password string) ([]string, error) {
	login, err := getUserWithUsername(username)
	if CheckError(err) {

		loginFail(w, r)
	} else {

		log.Println(login.Password)
		if login.Password == MD5(password) {
			fmt.Println("\ngood Password")
			log.Println("connecter avec : " + username)
			checkIfUserIsBanned(username)
			log.Println(data.Warning.ban)
			if data.Warning.ban != 1 {
				loginSuccess(w, r, login)
			} else {
				log.Println("Un banni essaie de ce connecter")
			}
		} else {
			fmt.Println("\nwrong Password")
			loginFail(w, r)
		}
	}

	return nil, nil
}

func loginSuccess(w http.ResponseWriter, r *http.Request, auth User) {
	data.User = auth
	data.Login = true
	log.Print(data.User.Email)
	Redirect(w, r, "/home")
}

func loginFail(w http.ResponseWriter, r *http.Request) {
	data.User = User{}
	data.Login = false
	Redirect(w, r, "/login")
}
