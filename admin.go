package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Banned struct {
	name string
	ban  int
	time string
}

func ban(name string, ban int, day int, month int, years int) error {
	quick, err := db.Query(fmt.Sprintf("INSERT INTO `ban` (`user_name`, `ban`, `time_ban`) VALUES('%s', '%d', '%s')", name, ban, time.Now().AddDate(day, month, years)))
	CheckError(err)
	defer quick.Close()
	return nil
}
func checkIfUserIsBanned(name string) (Banned, error) {
	checkBan := Banned{}
	lookBanTable := fmt.Sprintf("SELECT * FROM ban WHERE user_name='%s'", name)
	err := db.QueryRow(lookBanTable).Scan(&checkBan.name, &checkBan.ban, &checkBan.time)
	log.Println("Ban : ", checkBan.ban)
	data.Warning.ban = checkBan.ban
	return checkBan, err
}
func Hours(day int, month int, years int) {
	time.Now().AddDate(day, month, years)
}
func new_tag(name string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `tags` (`name`) VALUES ('%s')", name))
	log.Print("Le tag : [" + name + "] a bien été implémenté")
	CheckError(err)
	defer insert.Close()
	return nil
}
func HandleAdminPanel(w http.ResponseWriter, r *http.Request) {
	log.Print(data.User.PublicInfo.Admin)
	data.Page.Title = "Admin"
	data.Page.Style = "admin"
	pseudo := r.FormValue("pseudo")
	day := r.FormValue("day")
	month := r.FormValue("month")
	year := r.FormValue("years")
	tag := r.FormValue("tag")
	if data.User.PublicInfo.Admin == 1 {
		log.Println("go")
		if r.FormValue("envoyer") == "Envoyer" {
			log.Print("let's go !")
			new_tag(tag)
		}
	}
	if r.FormValue("good_bye") != "" {
		ban(pseudo, 1, Atoi(day), Atoi(month), Atoi(year))
	}
	if data.User.PublicInfo.Admin == 0 {
		go log.Print("L'utilisateur n'es pas un admin")
		Redirect(w, r, "/home")
	}
	tmpl.ExecuteTemplate(w, "admin", data)
}
