package main

import (
	"database/sql"
	"fmt"
	"log"

	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	Page       PageData
	Auth       Auth
	Error      string
	code       string
	Authorized bool
}

type PageData struct {
	Title string
	Style string
}

type Auth struct {
	user_id    string
	username   string
	password   string
	email      string
	phone      string
	first_name string
	last_name  string
	address    string
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

var tmpl *template.Template
var db *sql.DB
var data Data = Data{}

const (
	isValidEmail = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	normal_user  = 0
	Port         = "4443"
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
		data.Auth = Auth{}
		data.Error = ""
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

func handle404(w http.ResponseWriter, r *http.Request) {

	data.Page.Title = "Error 404"
	data.Page.Style = "error"

	errors := r.FormValue("click")
	if errors == "Accueil" {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
	}
	log.Print(errors)
	tmpl.ExecuteTemplate(w, "404", data)
}

func main() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/goeasydb")
	defer db.Close()

	fmt.Println("La connection a bien été établie")

	pageServer := http.FileServer(http.Dir("static/html"))
	http.Handle("/html/", http.StripPrefix("/html/", pageServer))

	styleServer := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styleServer))

	scriptServer := http.FileServer(http.Dir("static/js"))
	http.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	imageServer := http.FileServer(http.Dir("static/images"))
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))

	tmpl = template.Must(template.ParseGlob("static/html/*.html"))

	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/code", handleGetCode)
	http.HandleFunc("/", handle404)

	http.HandleFunc("/mangetesmort", googleLogin)
	http.HandleFunc("/google/callback", googleCallback)

	print("Lancement de la page instanciée sur : " + Host + ":" + Port + "/register")
	http.ListenAndServe(Host+":"+Port, nil)

}
