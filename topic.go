package main

import (
	"fmt"
	"net/http"
)

func createTopicInDB(name string, text_post string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `topic` (`name`, `user_id`, `contain`) VALUES ('%s','%s', '%s')", name, data.Auth.user_id, text_post))
	checkError(err)
	defer insert.Close()
	return nil
}
func handleGetTopic(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Topic"
	data.Page.Style = "topic"
	Name_Topic := r.FormValue("name_topic")
	contain := r.FormValue("contenue")
	if data.Login == true {
		if r.FormValue("submit") == "Envoyer" {
			if len(Name_Topic) >= 2 {
				if len(contain) > 0 {
					createTopicInDB(Name_Topic, contain)
				} else {
					var error_text = `<p style="color: red;">Le post doit contenire du text `
					fmt.Fprintf(w, error_text)
				}
			} else {
				var error_NameTopic = `<p style="color: red;">Le nom du topic est trop petit `
				fmt.Fprintf(w, error_NameTopic)
			}
		}
	} else {
		var error_connexion = `<p style="color: red;">La création de topic es réserver aux utilisateurs `
		fmt.Fprintf(w, error_connexion)
	}
	tmpl.ExecuteTemplate(w, "topic", data)
}
