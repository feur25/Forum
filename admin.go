package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Banned struct {
	name  string
	ban   int
	day   string
	month string
	year  string
}

func ban(name string, ban int, day string, month string, years string) error {
	quick, err := db.Query(fmt.Sprintf("INSERT INTO `ban` (`user_name`, `ban`, `day`, `month`, `year`) VALUES('%s', '%d', '%s', '%s', '%s')", name, ban, day, month, years))
	CheckError(err)
	defer quick.Close()
	return nil
}
func removeBan(name string) error {
	update, err := db.Query(fmt.Sprintf("DELETE FROM `ban` WHERE `user_name`='%s'", name))
	CheckError(err)
	defer update.Close()
	return nil
}
func checkIfUserIsBanned(name string) (Banned, error) {
	checkBan := Banned{}
	lookBanTable := fmt.Sprintf("SELECT * FROM ban WHERE user_name='%s'", name)
	err := db.QueryRow(lookBanTable).Scan(&checkBan.name, &checkBan.ban, &checkBan.day, &checkBan.month, &checkBan.year)
	//log.Println("Ban : ", checkBan.ban)
	time_today := time.Now().String()
	convertTimeTodayToInt := Atoi(time_today[:4] + time_today[5:7] + time_today[8:10])
	countBan := Atoi(checkBan.year + checkBan.month + checkBan.day)
	//log.Println(convertTimeTodayToInt, " + ", countBan)
	if convertTimeTodayToInt > countBan {
		removeBan(name)
	}
	data.Warning.ban = checkBan.ban
	return checkBan, err
}
func Hours(day int, month int, years int) {
	time.Now().AddDate(day, month, years)
}
func new_tag(name string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `tags` (`name`) VALUES ('%s')", name))
	//log.Print("Le tag : [" + name + "] a bien été implémenté")
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
	topicId := r.FormValue("topic_id")
	if data.User.PublicInfo.Admin == 1 {
		log.Println("go")
		if r.FormValue("envoyer") == "Envoyer" {
			new_tag(tag)
		}
		if r.FormValue("removetopic") != "" {
			removeTopicId(Atoi(topicId))
		}
	} else {
		go log.Print("L'utilisateur n'es pas un admin")
		Redirect(w, r, "/home")
	}
	if r.FormValue("good_bye") != "" {
		ban(pseudo, 1, day, month, year)
	}
	tmpl.ExecuteTemplate(w, "admin-page", data)
}
