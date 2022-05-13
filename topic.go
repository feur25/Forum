package main

import (
	"fmt"
	"net/http"
)

func getMostRecentTopics(length int) []TopicData {
	selection := fmt.Sprintf("SELECT topics.*, users.user_id, users.username, users.image_link, users.is_admin FROM topics LEFT JOIN users ON topics.user_id = users.user_id ORDER BY topics.creation_time ASC LIMIT %d", length)
	query, err := db.Query(selection)
	if checkError(err) {
		return nil
	}
	defer query.Close()

	result := []TopicData{}
	for query.Next() {
		var topic TopicData
		query.Scan(&topic.Topic.TopicId, &topic.Topic.CreationTime, &topic.Topic.Content, &topic.Topic.Name, &topic.Topic.UserId, &topic.Creator.Id, &topic.Creator.Username, &topic.Creator.ImageLink, &topic.Creator.Admin)
		ValidatePublicUserData(&topic.Creator)

		result = append(result, topic)
	}
	return result
}

func createTopicInDB(name string, PostText string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `topic` (`name`, `user_id`, `contain`) VALUES ('%s','%s', '%s')", name, data.User.PublicInfo.Id, PostText))
	checkError(err)
	defer insert.Close()
	return nil
}
func handleGetTopic(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Topic"
	data.Page.Style = "topic"
	Name_Topic := r.FormValue("name_topic")
	contain := r.FormValue("contenue")
	if data.Login {
		if r.FormValue("submit") == "Envoyer" {
			if len(Name_Topic) >= 2 {
				if len(contain) > 0 {
					createTopicInDB(Name_Topic, contain)
				} else {
					var error_text = `<p style="color: red;">Le post doit contenire du text `
					fmt.Fprint(w, error_text)
				}
			} else {
				var error_NameTopic = `<p style="color: red;">Le nom du topic est trop petit `
				fmt.Fprint(w, error_NameTopic)
			}
		}
	} else {
		var error_connexion = `<p style="color: red;">La création de topic es réserver aux utilisateurs `
		fmt.Fprint(w, error_connexion)
	}
	tmpl.ExecuteTemplate(w, "topic", data)
}
