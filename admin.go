package main

import (
	"fmt"
	"log"
	"net/http"
)

func new_tag(name string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `tags` (`name`) VALUES ('%s')", name))
	checkError(err)
	defer insert.Close()
	return nil
}
func check_if_admin() {
	insert, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE  is_admin = '%d'", data.Admin))
	checkError(err)
	defer insert.Close()
}
func handleAdminPanel(w http.ResponseWriter, r *http.Request) {
	check_if_admin()
	log.Print(data.Admin)
	data.Page.Title = "Admin"
	data.Page.Style = "admin"
	if data.Admin == 1 {
		new_tag(r.FormValue("tag"))
	} else {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
	}
	tmpl.ExecuteTemplate(w, "admin", data)
}
