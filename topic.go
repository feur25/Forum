package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TopicList struct {
	SelectedTopic TopicData
	Topics        []TopicData
}

type TopicData struct {
	Topic        Topic
	Creator      PublicUser
	Message      TopicMessage
	TopicMessage []TopicMessage
}
type Topic struct {
	TopicId      int
	CreationTime string
	Content      string
	Name         string
	UserId       string
}
type TopicMessage struct {
	CommentId int
	Date      string
	Content   string
	UserId    int
	TopicId   int
	UserName  string
}

func HandleTopic(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Topic"
	data.Page.Style = "topic"
	comment := r.FormValue("contenue")

	topicId, err := GetUrlParam(r, "id")
	if CheckError(err) {
		data.Page.TopicList.SelectedTopic = TopicData{}

	} else {
		data.Page.TopicList.SelectedTopic = getTopicWithId(Atoi(topicId))
		log.Println("ici : ", data.Page.TopicList.SelectedTopic.Message.Content)

	}
	if comment != "" && r.FormValue("submit") != "" {
		responseTopic(comment, Atoi(topicId))
	}
	log.Println(Atoi(topicId))
	data.Page.TopicList.Topics = getMostRecentTopics(20)
	data.Page.TopicList.SelectedTopic.Message.TopicId = Atoi(topicId)
	getFriendships(&data.Page.FriendList)
	data.Page.TopicList.Topics = getMostRecentMessageTopics(20)
	log.Println(getMostRecentMessageTopics(20))
	log.Println("la : ", data.Page.TopicList.SelectedTopic.Message.Content)

	tmpl.ExecuteTemplate(w, "topic", data)
}

func getTopicWithId(topicId int) TopicData {
	selection := fmt.Sprintf("SELECT topics.*, users.user_id, users.username, users.image_link, users.is_admin FROM topics LEFT JOIN users ON topics.user_id = users.user_id WHERE topics.topic_id = %d", topicId)
	topic := TopicData{}
	err := db.QueryRow(selection).Scan(&topic.Topic.TopicId, &topic.Topic.CreationTime, &topic.Topic.Content, &topic.Topic.Name, &topic.Topic.UserId, &topic.Creator.Id, &topic.Creator.Username, &topic.Creator.ImageLink, &topic.Creator.Admin)
	CheckError(err)
	ValidatePublicUserData(&topic.Creator)

	return topic
}

func getMostRecentTopics(length int) []TopicData {
	selection := fmt.Sprintf("SELECT topics.*, users.user_id, users.username, users.image_link, users.is_admin FROM topics LEFT JOIN users ON topics.user_id = users.user_id ORDER BY topics.creation_time ASC LIMIT %d", length)
	query, err := db.Query(selection)
	if CheckError(err) {
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

func createTopic(name string, PostText string) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `topic` (`name`, `user_id`, `contain`) VALUES ('%s','%s', '%s')", name, data.User.PublicInfo.Id, PostText))
	CheckError(err)
	defer insert.Close()
	return nil
}
func responseTopic(text string, id int) error {
	userId, _ := strconv.Atoi(data.User.PublicInfo.Id)
	comment, err := db.Query(fmt.Sprintf("INSERT INTO `comments` (`creation_date`, `content`, `user_id`, `topic_id`,`user_name`) VALUES ('%s','%s', '%d','%d','%s')", datetime, text, userId, id, data.User.PublicInfo.Username))
	CheckError(err)
	defer comment.Close()
	return nil
}
func getTopicMessageWithId(topicId int) TopicData {
	message := fmt.Sprintf("SELECT comments.*, users.user_id, users.username, users.image_link, users.is_admin FROM comments LEFT JOIN users ON comments.user_id = users.user_id WHERE comments.topic_id = %d", topicId)
	getMessage := TopicData{}
	err := db.QueryRow(message).Scan(&getMessage.Message.CommentId, &getMessage.Message.Date, &getMessage.Message.Content, &getMessage.Message.UserId, &getMessage.Message.TopicId, &getMessage.Message.UserName, &getMessage.Creator.Id, &getMessage.Creator.Username, &getMessage.Creator.ImageLink, &getMessage.Creator.Admin)
	CheckError(err)
	ValidatePublicUserData(&getMessage.Creator)

	return getMessage
}
func getMostRecentMessageTopics(length int) []TopicData {
	convert := strconv.Itoa(data.Page.TopicList.SelectedTopic.Message.TopicId)
	messageMostRecent := fmt.Sprintf("SELECT comments.*, users.user_id, users.username, users.image_link, users.is_admin FROM comments LEFT JOIN users ON comments.user_id = users.user_id WHERE comments.topic_id= '%s' ORDER BY comments.creation_date ASC LIMIT %d", convert, length)
	query, err := db.Query(messageMostRecent)
	if CheckError(err) {
		return nil
	}
	defer query.Close()

	result := []TopicData{}
	for query.Next() {
		getMessage := TopicData{}
		query.Scan(&getMessage.Message.CommentId, &getMessage.Message.Date, &getMessage.Message.Content, &getMessage.Message.UserId, &getMessage.Message.TopicId, &getMessage.Message.UserName, &getMessage.Creator.Id, &getMessage.Creator.Username, &getMessage.Creator.ImageLink, &getMessage.Creator.Admin)
		ValidatePublicUserData(&getMessage.Creator)
		log.Println(getMessage.Message.Content)

		result = append(result, getMessage)
	}
	return result
}

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
