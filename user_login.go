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
	log.Print("pseudo : " + username + " mdp : " + password)
	if isButtonPressed(r, "login") {
		checkUserLogin(w, r, username, password)
	} else if isButtonPressed(r, "disconnect") {
		disconnect(w, r)
	}
	tmpl.ExecuteTemplate(w, "login", data)
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

func disconnect(w http.ResponseWriter, r *http.Request) {
	data.User = User{}
	data.Login = false
	log.Print("wayouuu")
	http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
}

func loginSuccess(w http.ResponseWriter, r *http.Request, auth User) {
	data.User = auth
	data.Login = true
	log.Print(data.User.Email)
	http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
}

func loginFail(w http.ResponseWriter, r *http.Request) {
	data.User = User{}
	data.Error = "wrong password"
	http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
}
