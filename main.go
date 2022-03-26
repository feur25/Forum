package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	// "go/types"
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

func handleFunc(w http.ResponseWriter, r *http.Request) {

	tmpl.Execute(w, data)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func encrypt(plainText string) string {
	encryptedText := sha256.Sum256([]byte(plainText))
	return string(encryptedText[:])
}

func createUser(username string, email string, phone string, firstName string, lastName string, address string, date string, password string, admin int) error {
	var datetime = time.Now().UTC().Format("2006-01-02 03:04:05")
	data.time = datetime
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `utilisateur`(`username`, `email`, `phone_number`, `first_name`, `last_name`, `address`, `creation_date`, `password`, `is_admin`) VALUES ('%s','%s','%s','%s','%s','%s', '%s', '%s', '%d')", username, email, phone, firstName, lastName, address, datetime, password, admin))
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	fmt.Print("Mission réusise")

	return nil
}

func userLogin(username, password string) ([]string, error) {
	selectQuery, err := db.Query(fmt.Sprintf("SELECT * FROM `utilisateur` WHERE username = '%s' OR email = '%s' LIMIT 1", username, username))
	checkError(err)

	columns, err := selectQuery.Columns()
	checkError(err)
	fmt.Println(columns)

	if columns[9] != password {
		return []string{}, errors.New("wrong Password")
	}

	selectQuery.Close()

	return columns, nil
}

func checkPasswordValidity(password string) bool {

	if len(data.password) < 10 {
		log.Print("Le mot de passe doit contenir 10 caractères ou plus.")
		return false
	}

	upper, err := regexp.MatchString(data.password, containsUpperCase)
	checkError(err)
	if !upper {
		log.Print("Le mot de passe ne contient pas de lettres majuscules.")
		return false
	}

	lower, err := regexp.MatchString(data.password, containsLowerCase)
	checkError(err)
	if !lower {
		log.Print("Le mot de passe ne contient pas de lettres minuscule.")
		return false
	}

	number, err := regexp.MatchString(data.password, containsNumber)
	checkError(err)
	if !number {
		log.Print("Le mot de passe ne contient pas de chiffres.")
		return false
	}

	special, err := regexp.MatchString(data.password, containsSpecial)
	checkError(err)
	if !special {
		log.Print("Le mot de passe ne contient pas de caractères spéciaux.")
		return false
	}

	return true
}
func checkemail(email string) bool {
	if !strings.Contains(data.email, "@") {
		return false
	}
	if !strings.Contains(data.email, ".") {
		return false
	}
	return true
}
func user_create(w http.ResponseWriter, r *http.Request) {
	register := template.Must(template.ParseFiles("static/html/register.html"))
	data.username = r.FormValue("username")
	data.password = r.FormValue("password")
	data.adress = r.FormValue("address")
	data.phone = r.FormValue("phone")
	data.email = r.FormValue("email")
	data.first_name = r.FormValue("first_name")
	data.last_name = r.FormValue("last_name")
	if len(data.username) >= 4 && len(data.adress) != 0 && len(data.phone) >= 10 && checkemail(data.email) && len(data.first_name) >= 4 && len(data.last_name) >= 3 {
		createUser(data.username, data.email, data.phone, data.first_name, data.last_name, data.adress, data.time, data.password, normal_user)
		http.Redirect(w, r, "http://"+Host+":"+Port+"/login", 301)
	}
	register.Execute(w, data)
}
func Handle404(w http.ResponseWriter, r *http.Request) {
	Error := template.Must(template.ParseFiles("static/html/error.html"))
	errors := r.FormValue("click")
	if errors == "Accueil" {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", 301)
	}
	log.Print(errors)
	Error.Execute(w, r)
}
func userAuth(w http.ResponseWriter, r *http.Request) {
	login := template.Must(template.ParseFiles("static/html/login.html"))
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Print("\n" + username)
	log.Print(password)
	envoyer := r.FormValue("submit")
	if envoyer == "Envoyer" {
		userLogin(username, password)
	}
	login.Execute(w, data)
}

func main() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test_forum")
	fmt.Println("Go Mysql Tutorial")
	defer db.Close()

	fmt.Println("La connection a bien été établie")

	fmt.Print("Nice")

	pageServer := http.FileServer(http.Dir("static/html"))
	http.Handle("/html/", http.StripPrefix("/html/", pageServer))
	http.HandleFunc("/register", user_create)
	http.HandleFunc("/login", userAuth)
	http.HandleFunc("/", Handle404)

	styleServer := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styleServer))

	scriptServer := http.FileServer(http.Dir("static/js"))
	http.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	imageServer := http.FileServer(http.Dir("static/images"))
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))
	print("Lancement de la page instancier sur : " + Host + ":" + Port + "/register")
	http.ListenAndServe(Host+":"+Port, nil)

	http.HandleFunc("/", handleFunc)
}
