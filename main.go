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
	user_id    string
	username   string
	password   string
	adress     string
	phone      string
	email      string
	first_name string
	last_name  string
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

var tmpl = template.Must(template.ParseGlob("static/html/*.html"))
var db *sql.DB
var data Auth = Auth{}

const (
	isValidEmail = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	normal_user  = 0
	Port         = "4444"
	Host         = "localhost"
)

func checkErrorPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func checkErrorLogout(err error) bool {
	if err != nil {
		log.Println(err.Error())
		data = Auth{}
		return true
	}
	return false
}
func checkError(err error) bool {
	if err != nil {
		log.Println(err.Error())
		return true
	}
	return false
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	register := template.Must(template.ParseFiles("static/html/register.html"))

	username := r.FormValue("username")
	password := r.FormValue("password")
	address := r.FormValue("address")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	submit := r.FormValue("submit")
	if submit != "" {

		usernameValidity, err := checkUsernameValidity(username)
		checkErrorLogout(err)
		addressValidity, err := checkAdressValidity(address)
		checkErrorLogout(err)
		emailValidity, err := checkEmailValidity(email)
		checkErrorLogout(err)
		passValidity, err := checkPasswordValidity(w, r, password)
		checkErrorLogout(err)
		fmt.Print(w, err.Error())

		if usernameValidity && addressValidity && len(phone) >= 10 && emailValidity && passValidity && len(first_name) >= 4 && len(last_name) >= 3 {

			err := createUserInDatabase(username, email, phone, first_name, last_name, address, password, normal_user)
			if !checkErrorLogout(err) {
				http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
			} else {
				http.Redirect(w, r, "http://"+Host+":"+Port+"/register", http.StatusMovedPermanently)
			}
		}

	} else {
		data = Auth{}
	}
	submit = ""
	register.Execute(w, data)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	login := template.Must(template.ParseFiles("static/html/login.html"))
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Print("\n" + username)
	log.Print(password)
	if r.FormValue("submit") == "Envoyer" {
		checkUserLogin(w, r, username, password)
	}
	login.Execute(w, data)
}

func settingUser(w http.ResponseWriter, r *http.Request) {
	update := template.Must(template.ParseFiles("static/html/update.html"))
	password := r.FormValue("pasword")
	passwordcheck := r.FormValue("checkpassword")
	log.Print("\n" + passwordcheck)
	log.Print(password)
	if r.FormValue("submit") == "Envoyer" {
		if password == passwordcheck {
			updateUserInDatabase(password)
		}
	}
	update.Execute(w, data)
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
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/goeasydb")
	defer db.Close()

	fmt.Println("La connection a bien été établie")

	pageServer := http.FileServer(http.Dir("static/html"))
	http.Handle("/html/", http.StripPrefix("/html/", pageServer))
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/", handle404)

	http.HandleFunc("/mangetesmort", googleLogin)
	http.HandleFunc("/google/callback", googleCallback)

	styleServer := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styleServer))

	scriptServer := http.FileServer(http.Dir("static/js"))
	http.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	imageServer := http.FileServer(http.Dir("static/images"))
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))
	print("Lancement de la page instancier sur : " + Host + ":" + Port + "/register")
	http.ListenAndServe(Host+":"+Port, nil)

}
