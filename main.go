package main

import (
	"database/sql"
	"fmt"
	"log"

	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Auth struct {
	username   string
	password   string
	time       string
	adress     string
	phone      string
	email      string
	first_name string
	last_name  string
}

var tmpl = template.Must(template.ParseGlob("static/html/*.html"))
var db *sql.DB
var data Auth = Auth{}

const (
	containsUpperCase = "^[A-Z]$"
	containsLowerCase = "^[a-z]$"
	containsNumber    = "^[0-9]$"
	containsSpecial   = `^[-+_!@#$%^&*.,?\/\\]$`
	normal_user       = 0
	Port              = "4448"
	Host              = "localhost"
)

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	register := template.Must(template.ParseFiles("static/html/register.html"))
	data.username = r.FormValue("username")
	data.password = r.FormValue("password")
	data.adress = r.FormValue("address")
	data.phone = r.FormValue("phone")
	data.email = r.FormValue("email")
	data.first_name = r.FormValue("first_name")
	data.last_name = r.FormValue("last_name")

	emailValidity, err := checkEmailValidity(data.email)
	checkError(err)
	passValidity, err := checkPasswordValidity(data.password)
	checkError(err)
	if len(data.username) >= 4 && len(data.adress) != 0 && len(data.phone) >= 10 && emailValidity && emailemailValidity && len(data.first_name) >= 4 && len(data.last_name) >= 3 {

		createUserInDatabase(data.username, data.email, data.phone, data.first_name, data.last_name, data.adress, data.time, data.password, normal_user)
		http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
	}
	register.Execute(w, data)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	login := template.Must(template.ParseFiles("static/html/login.html"))
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Print("\n" + username)
	log.Print(password)
	envoyer := r.FormValue("submit")
	if envoyer == "Envoyer" {
		CheckuserLogin(username, password)
	}
	login.Execute(w, data)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	Error := template.Must(template.ParseFiles("static/html/error.html"))
	errors := r.FormValue("click")
	if errors == "Accueil" {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
	}
	log.Print(errors)
	Error.Execute(w, r)
}

func main() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test_forum")
	defer db.Close()

	fmt.Println("La connection a bien été établie")

	pageServer := http.FileServer(http.Dir("static/html"))
	http.Handle("/html/", http.StripPrefix("/html/", pageServer))
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/", handle404)

	styleServer := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styleServer))

	scriptServer := http.FileServer(http.Dir("static/js"))
	http.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	imageServer := http.FileServer(http.Dir("static/images"))
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))
	print("Lancement de la page instancier sur : " + Host + ":" + Port + "/register")
	http.ListenAndServe(Host+":"+Port, nil)

}
