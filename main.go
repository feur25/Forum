package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	// "go/types"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Auth struct {
	username string
	password string
	time     string
}

var tmpl = template.Must(template.ParseGlob("static/html/*.html"))
var db *sql.DB
var data Auth = Auth{}

const (
	containsUpperCase = "^[A-Z]$"
	containsLowerCase = "^[a-z]$"
	containsNumber    = "^[0-9]$"
	containsSpecial   = `^[-+_!@#$%^&*.,?\/\\]$`
)

func userAuth(w http.ResponseWriter, r *http.Request) {
	data.username = r.FormValue("username")
	data.password = r.FormValue("password")

	userLogin(data.username, data.password)
}

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

/*('127','AxelSeven','axelsevenet@gmail.com',true,'"0616694403"','Axel','Sevenet','6 Butte des 3 Moulins','JesuisDieu05!',25/03/2022)*/
func createUser(id string, username string, email string, phone string, firstName string, lastName string, address string, date string, password string, admin int) error {
	var datetime = time.Now().UTC().Format("2006-01-02 03:04:05")
	data.time = datetime
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `utilisateur`(`user_id`, `username`, `email`, `phone_number`, `first_name`, `last_name`, `address`, `creation_date`, `password`, `is_admin`) VALUES ('%s','%s','%s','%s','%s','%s','%s', '%s', '%s', '%d')", id, username, email, phone, firstName, lastName, address, datetime, password, admin))
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	fmt.Print("Mission réusise")

	return nil
}

func userLogin(username, password string) ([]string, error) {
	selectQuery, err := db.Query(fmt.Sprintf("SELECT * INTO `utilisateur` WHERE username = `%s` OR email = `%s` LIMIT 1", username, username))
	checkError(err)

	columns, err := selectQuery.Columns()
	checkError(err)
	fmt.Println(columns)

	if columns[3] != encrypt(password) {
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

func main() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test_forum")
	fmt.Println("Go Mysql Tutorial")
	defer db.Close()
	createUser("126", "AxelSeven5", "axelsevenet5@gmail.com", "0616694403", "Axel", "Sevenet", "6 Butte des 3 Moulins", data.time, "mabite", 3)
	/*userLogin("AxelSeven3", "mabite")
	userLogin("axelsevenet3@gmail.com", "mabite")*/

	fmt.Println("La connection a bien été établie")

	fmt.Print("Nice")

	pageServer := http.FileServer(http.Dir("static/html"))
	http.Handle("/html/", http.StripPrefix("/html/", pageServer))

	styleServer := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styleServer))

	scriptServer := http.FileServer(http.Dir("static/js"))
	http.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	imageServer := http.FileServer(http.Dir("static/images"))
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))

	http.HandleFunc("/", handleFunc)
}
