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

func createUser(username, email, password, phone, firstName, lastName, address string, admin bool) error {
	if !checkPasswordValidity(password) {
		return errors.New("password not valid")
	}
	dt := time.Now()
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `users`(`username`, `email`, `password`, `phone_number`, `first_name`, `last_name`, `address`, `is_admin`, `creation_data`) VALUES ('%s','%s', %s,'%s','%s','%s','%s','%t', `%s`)", username, email, encrypt(password), phone, firstName, lastName, address, admin, dt.Format("yyyy-mm-dd")))
	checkError(err)
	insert.Close()

	return nil
}
func userLogin(username, password string) ([]string, error) {
	selectQuery, err := db.Query(fmt.Sprintf("SELECT * INTO `users` WHERE username = `%s` OR email = `%s` LIMIT 1", username, username))
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
	fmt.Println("Go Mysql Tutorial")
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_easy")
	checkError(err)
	defer db.Close()

	createUser("AxelSeven", "axelsevenet@gmail.com", "mabite", "0616694403", "Axel", "Sevenet", "6 Butte des 3 Moulins", true)
	userLogin("AxelSeven", "mabite")
	userLogin("axelsevenet@gmail.com", "mabite")

	fmt.Println("La connection a bien été établie")

	insert, err := db.Query("INSERT INTO `products`(`productCode`, `productName`, `productLine`, `productScale`, `productVendor`, `productDescription`, `quantityInStock`, `buyPrice`, `MSRP`) VALUES ('S10_1100','humberger','Planes','1:18','Motor City Art Classics ','American supersonic, twin-engine, two-seat, twin-tail, variable-sweep wing fighter aircraft ',200,88.20,175)")
	checkError(err)
	insert.Close()

	fmt.Print("Nice")

	pageServer := http.FileServer(http.Dir("static/html"))
	http.Handle("/html/", http.StripPrefix("/html/", pageServer))

	styleServer := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styleServer))

	scriptServer := http.FileServer(http.Dir("static/js"))
	http.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	imageServer := http.FileServer(http.Dir("static/images"))
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))

	tmpl = template.Must(template.ParseGlob("html/*.html"))

	http.HandleFunc("/", handleFunc)
}
