package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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
	friend_info            Friend
	button                 bool
	Error                  string
	UpdateConfirmationCode string
	DeleteConfirmationCode string
	Authorized             bool
	Login                  bool
	Admin                  int
}
type session struct {
	expiry time.Time
}

type PageData struct {
	Topics  []TopicData
	Friends FriendList
	Title   string
	Style   string
}

type FriendList struct {
	AcceptedRequests []Friend
	SentRequests     []Friend
	ReceivedRequests []Friend
}

type Friend struct {
	FriendName  string
	FriendId    string
	SenderId    string
	RecipientId string
	Pending     int
}

type TopicData struct {
	Topic   Topic
	Creator PublicUser
}
type Topic struct {
	TopicId      int
	CreationTime string
	Content      string
	Name         string
	UserId       string
}

var tmpl *template.Template
var db *sql.DB
var data Data = Data{}

const (
	isValidEmail = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	normal_user  = 0
	Port         = "4448"
	Host         = "localhost"
)

func MD5(raw string) string {
	encrypted := md5.Sum([]byte(raw))
	return hex.EncodeToString(encrypted[:])
}

func checkErrorLogout(err error) bool {
	if err != nil {
		log.Println(err.Error())
		data.Page.Topics = []TopicData{}
		data.User = User{}
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
func isButtonPressed(r *http.Request, buttonName string) bool {
	return r.FormValue(buttonName) != ""
}
func convertIntToString(text string) string {
	strconv.ParseInt(text, 10, 64)
	return text
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

	tmpl = template.Must(template.ParseGlob("static/html/*.html"))

	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	// http.HandleFunc("/code", handleGetCode)
	http.HandleFunc("/topic/", handleTopic)
	http.HandleFunc("/admin", handleAdminPanel)
	http.HandleFunc("/update", handleUpdateUser)
	http.HandleFunc("/delete", handleDeleteUser)
	http.HandleFunc("/friendreq", HandleSendFriendRequest)
	http.HandleFunc("/home", handleHome)
	http.HandleFunc("/", handle404)

	http.HandleFunc("/mangetesmort", googleLogin)
	http.HandleFunc("/google/callback", googleCallback)

	print("Lancement de la page effectué : " + Host + ":" + Port + "/register")
	http.ListenAndServe(Host+":"+Port, nil)
}
