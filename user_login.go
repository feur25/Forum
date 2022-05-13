package main

import (
	"fmt"
	"log"

	"net/http"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {

	data.Page.Title = "Login"
	data.Page.Style = "login"

	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Print("\n" + username)
	log.Print(password)
	if r.FormValue("submit") == "Envoyer" {
		checkUserLogin(w, r, username, password)
	}
	tmpl.ExecuteTemplate(w, "login", data)
}

func loginSuccess(w http.ResponseWriter, r *http.Request, auth User) {
	data.User = auth
	data.User.PublicInfo.Username = r.FormValue("username")
	data.Login = true
	log.Print(data.User.Email)
	http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
}

func loginFail(w http.ResponseWriter, r *http.Request) {
	data.User = User{}
	data.Error = "wrong password"
	http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
}

func checkUserLogin(w http.ResponseWriter, r *http.Request, username, password string) ([]string, error) {
	login, err := getUserInDB(username)
	if checkError(err) {

		loginFail(w, r)
	} else {

		log.Println(login.Password)
		if login.Password == MD5(password) {
			fmt.Println("\ngood Password")
			loginSuccess(w, r, login)
		} else {
			fmt.Println("\nwrong Password")
			loginFail(w, r)
		}
	}

	return nil, nil
}
