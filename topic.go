package main

import (
	"fmt"
	"log"
	"net/http"
)

type TopicList struct {
	SelectedTopic Topic
	Topics        []Topic
}

type Topic struct {
	TopicId      int
	CreationTime string
	Name         string
	Content      string
	Creator      PublicUser
	Comments     []TopicComment
}
type TopicComment struct {
	CommentId    int
	TopicId      int
	CreationTime string
	Content      string
	Creator      PublicUser
}

func HandleTopic(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Topic"
	data.Page.Style = "topic"

	topicId, err := GetUrlParam(r, "id")
	if IsButtonPressed(r, "submit") {
		comment := r.FormValue("comment")
		if data.Login {
			commentTopic(topicId, comment)
		} else {
			log.Println("Il faut vous connecter !!")
		}
	}

	if CheckError(err) {
		data.Page.TopicList.SelectedTopic = Topic{}
		data.Page.TopicList.SelectedTopic.Comments = []TopicComment{}
	} else {
		data.Page.TopicList.SelectedTopic = getTopicWithId(Atoi(topicId))
		data.Page.TopicList.SelectedTopic.Comments = getTopicComments(topicId, 20)
	}

	tmpl.ExecuteTemplate(w, "topic-page", data)
}
func HandlecreateTopic(w http.ResponseWriter, r *http.Request) {
	name_topic := r.FormValue("nametopic")
	text_topic := r.FormValue("texttopic")
	log.Println("exemple")
	if data.Login {
		log.Println("exemple")
		if r.FormValue("post") == "Submit" {
			log.Println("exemple")
			createTopic(name_topic, text_topic)
		}
	}
	log.Println("exemple")
	tmpl.ExecuteTemplate(w, "topic-create-page", nil)
}

func getTopicWithId(topicId int) Topic {
	selection := fmt.Sprintf("SELECT topics.topic_id, topics.creation_time, topics.content, topics.name, users.user_id, users.username, users.image_link, users.is_admin FROM topics LEFT JOIN users ON topics.user_id = users.user_id WHERE topics.topic_id = %d", topicId)
	topic := Topic{}
	err := db.QueryRow(selection).Scan(&topic.TopicId, &topic.CreationTime, &topic.Content, &topic.Name, &topic.Creator.Id, &topic.Creator.Username, &topic.Creator.ImageLink, &topic.Creator.Admin)
	CheckError(err)

	selection = fmt.Sprintf("SELECT comments.comment_id, comments.topic_id, comments.creation_time, comments.content, users.user_id, users.username, users.image_link, users.is_admin FROM comments LEFT JOIN users ON comments.user_id = users.user_id WHERE comments.topic_id = %d", topicId)
	query, err := db.Query(selection)
	CheckError(err)
	defer query.Close()
	for query.Next() {
		var comment TopicComment
		query.Scan(&comment.CommentId, &comment.TopicId, &comment.CreationTime, &comment.Content, &comment.Creator.Id, &comment.Creator.Username, &comment.Creator.ImageLink, &comment.Creator.Admin)
		topic.Comments = append(topic.Comments, comment)
	}

	return topic
}
func removeTopicId(id int) error {
	rmTopic, err := db.Query("DELETE FROM `topics` WHERE topic_id='%d'", id)
	CheckError(err)
	defer rmTopic.Close()
	return nil
}
func getMostRecentTopics(length int) []Topic {
	selection := fmt.Sprintf("SELECT topics.topic_id, topics.creation_time, topics.content, topics.name, users.user_id, users.username, users.image_link, users.is_admin FROM topics LEFT JOIN users ON topics.user_id = users.user_id ORDER BY topics.creation_time ASC LIMIT %d", length)
	query, err := db.Query(selection)
	if CheckError(err) {
		return nil
	}
	defer query.Close()

	result := []Topic{}
	for query.Next() {
		var topic Topic
		query.Scan(&topic.TopicId, &topic.CreationTime, &topic.Content, &topic.Name, &topic.Creator.Id, &topic.Creator.Username, &topic.Creator.ImageLink, &topic.Creator.Admin)

		result = append(result, topic)
	}
	return result
}

func getTopicComments(topicId string, length int) []TopicComment {
	messageMostRecent := fmt.Sprintf("SELECT comments.comment_id, comments.creation_time, comments.content, comments.topic_id, users.user_id, users.username, users.image_link, users.is_admin FROM comments LEFT JOIN users ON comments.user_id = users.user_id WHERE comments.topic_id= '%s' ORDER BY comments.creation_time ASC LIMIT %d", topicId, length)
	query, err := db.Query(messageMostRecent)
	if CheckError(err) {
		return nil
	}
	defer query.Close()

	result := []TopicComment{}
	for query.Next() {
		comment := TopicComment{}
		query.Scan(&comment.CommentId, &comment.CreationTime, &comment.Content, &comment.TopicId, &comment.Creator.Id, &comment.Creator.Username, &comment.Creator.ImageLink, &comment.Creator.Admin)

		result = append(result, comment)
	}
	return result
}

func createTopic(name string, PostText string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `topics` (`content`, `name`, `user_id`) VALUES ('%s','%s', '%s')", name, PostText, data.User.PublicInfo.Id))
	CheckError(err)
	defer insert.Close()
	return nil
}
func commentTopic(topicId, commentText string) {
	insertion, err := db.Query(fmt.Sprintf("INSERT INTO `comments` (`content`, `user_id`, `topic_id`) VALUES ('%s', '%s', '%s')", commentText, data.User.PublicInfo.Id, topicId))
	CheckError(err)
	defer insertion.Close()
}

// func getTopicCommentWithId(topicId int) TopicData {
// 	message := fmt.Sprintf("SELECT comments.*, users.user_id, users.username, users.image_link, users.is_admin FROM comments LEFT JOIN users ON comments.user_id = users.user_id WHERE comments.topic_id = %d", topicId)
// 	getMessage := TopicData{}
// 	err := db.QueryRow(message).Scan(&getMessage.Message.CommentId, &getMessage.Message.Date, &getMessage.Message.Content, &getMessage.Message.UserId, &getMessage.Message.TopicId, &getMessage.Message.UserName, &getMessage.Creator.Id, &getMessage.Creator.Username, &getMessage.Creator.ImageLink, &getMessage.Creator.Admin)
// 	CheckError(err)
// 	// ValidatePublicUserData(&getMessage.Creator)

// 	return getMessage
// }

// func HandleGetTopic(w http.ResponseWriter, r *http.Request) {
// 	data.Page.Title = "Topic"
// 	data.Page.Style = "topic"
// 	Name_Topic := r.FormValue("name_topic")
// 	contain := r.FormValue("contenue")
// 	if data.Login {
// 		if r.FormValue("submit") == "Envoyer" {
// 			if len(Name_Topic) >= 2 {
// 				if len(contain) > 0 {
// 					createTopic(Name_Topic, contain)
// 				} else {
// 					var error_text = `<p style="color: red;">Le post doit contenire du text `
// 					fmt.Fprint(w, error_text)
// 				}
// 			} else {
// 				var error_NameTopic = `<p style="color: red;">Le nom du topic est trop petit `
// 				fmt.Fprint(w, error_NameTopic)
// 			}
// 		}
// 	} else {
// 		var error_connexion = `<p style="color: red;">La création de topic es réserver aux utilisateurs `
// 		fmt.Fprint(w, error_connexion)
// 	}
// 	tmpl.ExecuteTemplate(w, "topic", data)
// }
