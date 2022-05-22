package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	Page                   PageData
	User                   User
	Message                MessagePrivate
	Warning                Banned
	Error                  string
	UpdateConfirmationCode string
	DeleteConfirmationCode string
	Authorized             bool
	Login                  bool
	Admin                  int
}

type Key struct {
	keyPrivate string
	keyPublic  string
}

type PageData struct {
	TopicList  TopicList
	FriendList FriendList
	Title      string
	Style      string
}

var tmpl *template.Template
var db *sql.DB
var data Data = Data{}

const (
	isValidEmail = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	Port         = "4448"
	Host         = "localhost"
)

var datetime = time.Now().UTC().Format("2006-01-02 03:04:05")

func IsButtonPressed(r *http.Request, buttonName string) bool {
	return r.FormValue(buttonName) != ""
}

func Redirect(w http.ResponseWriter, r *http.Request, link string) {
	http.Redirect(w, r, "http://"+Host+":"+Port+link, http.StatusMovedPermanently)
}

func GetUrlParam(r *http.Request, param string) (string, error) {
	params := r.URL.Query()[param]
	if len(params) == 0 {
		return "", errors.New("parameter cannot be empty")
	}
	return params[0], nil
}

func MD5(raw string) string {
	encrypted := md5.Sum([]byte(raw))
	return hex.EncodeToString(encrypted[:])
}

func CheckError(err error) bool {
	if err != nil {
		log.Println(err.Error())
		return true
	}
	return false
}
func Atoi(x string) int {
	num, err := strconv.Atoi(x)
	if CheckError(err) {
		return 0
	}
	return num
}

func Handle404(w http.ResponseWriter, r *http.Request) {

	data.Page.Title = "Error 404"
	data.Page.Style = "error"

	errors := r.FormValue("click")
	if errors == "Accueil" {
		Redirect(w, r, "/home")
	}
	log.Print(errors)
	tmpl.ExecuteTemplate(w, "404", data)
}

func HttpHandle(url string, function func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		data.Error = ""

		function(w, r)
	})
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/goeasydb")
	defer db.Close()

	fmt.Println("La connection a bien été établie")

	pageServer := http.FileServer(http.Dir("static/html"))
	http.Handle("/html/", http.StripPrefix("/html/", pageServer))

	styleServer := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styleServer))

	scriptServer := http.FileServer(http.Dir("static/js"))
	http.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	imageServer := http.FileServer(http.Dir("static/img"))
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))

	var err error
	tmpl, err = template.New("").Funcs(template.FuncMap{
		"Translate": Translate,
		// "AcceptFriendRequest": AcceptFriendRequest,
		// "DenyFriendRequest":   DenyFriendRequest,
	}).ParseGlob("static/html/*.html")

	// tmpl, err = template.New("").ParseGlob("static/html/*.html")
	CheckError(err)

	// tmpl, err := tmpl.Funcs(template.FuncMap{
	// 	"AcceptFriendRequest": AcceptFriendRequest,
	// 	"DenyFriendRequest": DenyFriendRequest,
	// }).ParseGlob("static/html/*.html")
	// CheckError(err)

	HttpHandle("/register", HandleRegister)
	HttpHandle("/login", HandleLogin)
	// HttpHandle("/code", HandleGetCode)

	HttpHandle("/topic/", HandleTopic)
	HttpHandle("/admin", HandleAdminPanel)
	HttpHandle("/update", HandleUpdateUser)
	HttpHandle("/delete", HandleDeleteUser)
	HttpHandle("/cookie", indexHandler)
	HttpHandle("/home", HandleHome)
	HttpHandle("/friends/request/", HandleSendFriendRequest)
	HttpHandle("/friends/accept/", HandleAcceptFriendRequest)
	HttpHandle("/friends/deny/", HandleDenyFriendRequest)
	HttpHandle("/friends/delete/", HandleDeleteFriend)
	HttpHandle("/friends/message/", HandleMessageFriend)
	HttpHandle("/", Handle404)

	HttpHandle("/mangetesmort", googleLogin)
	HttpHandle("/google/callback", googleCallback)

	print("Lancement de la page effectué : " + Host + ":" + Port + "/register")
	http.ListenAndServe(Host+":"+Port, nil)
}
