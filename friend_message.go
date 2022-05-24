package main

import (
	"fmt"
	"log"
	"net/http"
)

type FriendMessage struct {
	Id               string
	RequestId        string
	Sender           PublicUser
	Recipient        PublicUser
	SenderContent    string
	RecipientContent string
	CreationTime     string
	UserIsSender     bool
}

func HandleMessageFriend(w http.ResponseWriter, r *http.Request) {

	requestId, err := GetUrlParam(r, "id")
	// if !checkFriendAuthorization(requestId, data.User.PublicInfo.Id) {
	// 	Redirect(w, r, "/home")
	// 	return
	// }

	if checkFriendAuthorization(requestId, data.User.PublicInfo.Id) && !CheckError(err) {

		data.Page.FriendList.SelectedFriend = getFriendRequest(requestId)
		data.Page.FriendList.SelectedFriendMessages = getFriendMessages(requestId)
	} else {
		data.Page.FriendList.SelectedFriend = Friend{}
		data.Page.FriendList.SelectedFriendMessages = []FriendMessage{}

		Redirect(w, r, "/home")
	}

	data.Page.Title = "Topic"
	data.Page.Style = "topic"

	tmpl.ExecuteTemplate(w, "friend-message-page", data)
}

func HandleSendMessage(w http.ResponseWriter, r *http.Request) {

	messageText, err := GetUrlParam(r, "message")
	if !CheckError(err) {
		message := insertFriendMessage(data.Page.FriendList.SelectedFriend.RequestId, data.User.PublicInfo.Id, data.Page.FriendList.SelectedFriend.FriendUser.Id, messageText, messageText)
		log.Print(message)
	}
	Redirect(w, r, "/friends/message?id="+data.Page.FriendList.SelectedFriend.RequestId)
}

func getFriendMessages(friendId string) []FriendMessage {
	var messages []FriendMessage
	selection := fmt.Sprintf("SELECT messages.message_id, messages.friend_id, sender.user_id, sender.username, sender.image_link, sender.is_admin, recipient.user_id, recipient.username, recipient.image_link, recipient.is_admin, messages.sender_content, messages.recipient_content, messages.creation_time FROM messages LEFT JOIN users as sender ON sender.user_id = messages.sender_id LEFT JOIN users as recipient ON recipient.user_id = messages.recipient_id WHERE messages.friend_id = '%s' ORDER BY messages.creation_time DESC", friendId)
	query, err := db.Query(selection)
	CheckError(err)
	defer query.Close()
	for query.Next() {
		message := FriendMessage{}
		err = query.Scan(&message.Id, &message.RequestId, &message.Sender.Id, &message.Sender.Username, &message.Sender.ImageLink, &message.Sender.Admin, &message.Recipient.Id, &message.Recipient.Username, &message.Recipient.ImageLink, &message.Recipient.Admin, &message.SenderContent, &message.RecipientContent, &message.CreationTime)
		CheckError(err)
		message.UserIsSender = message.Sender.Id == data.User.PublicInfo.Id
		messages = append(messages, message)
	}
	return messages
}

// func getFriendMessages(friendId string) []FriendMessage {
// 	var messages []FriendMessage
// 	selection := fmt.Sprintf("SELECT messages.message_id FROM messages WHERE messages.friend_id = '%s'", friendId)
// 	query, err := db.Query(selection)
// 	CheckError(err)
// 	defer query.Close()
// 	for query.Next() {
// 		message := FriendMessage{}
// 		query.Scan(&message.Id)
// 		log.Print(message)
// 		messages = append(messages, message)
// 	}
// 	return messages
// }

func insertFriendMessage(friendId, senderId, recipientId, senderContent, recipientContent string) FriendMessage {
	log.Print(friendId + " : " + senderId + " : " + recipientId + " : " + senderContent + " : " + recipientContent)
	insertion := fmt.Sprintf("INSERT INTO `messages` (`friend_id`, `sender_id`, `recipient_id`, `sender_content`, `recipient_content`) VALUES ( '%s', '%s', '%s','%s', '%s' )", friendId, senderId, recipientId, senderContent, recipientContent)
	_, err := db.Exec(insertion)
	CheckError(err)

	message := FriendMessage{}
	selection := fmt.Sprintf("SELECT messages.message_id, messages.friend_id, messages.friend_id, sender.user_id, sender.username, sender.image_link, sender.is_admin, recipient.user_id, recipient.username, recipient.image_link, recipient.is_admin, messages.sender_content, messages.recipient_content, messages.creation_time FROM friends LEFT JOIN users as sender ON sender.user_id = messages.sender_id LEFT JOIN users as recipient ON recipient.user_id = messages.recipient_id WHERE messages.friend_id = '%s' AND sender.user_id = '%s' AND recipient.user_id = '%s' AND messages.sender_content = '%s' AND messages.recipient_content = '%s' ", friendId, senderId, recipientId, senderContent, recipientContent)
	db.QueryRow(selection).Scan(&message.Id, &message.RequestId, &message.Sender.Id, &message.Sender.Username, &message.Sender.ImageLink, &message.Sender.Admin, &message.Recipient.Id, &message.Recipient.Username, &message.Recipient.ImageLink, &message.Recipient.Admin, &message.CreationTime)
	return message
}

// func createFriendMessage(senderId int, recipientName string, message string) error {
// 	var recipientId int
// 	var recipientEmail, senderName string
// 	selection := fmt.Sprintf("SELECT recipient.user_id, recipient.email, sender.username FROM users as recipient, users as sender WHERE sender.user_id = %d AND recipient.username = '%s'", senderId, recipientName)
// 	db.QueryRow(selection).Scan(&recipientId, &recipientEmail, &senderName)

// 	foundId := 0
// 	selection = fmt.Sprintf("SELECT `request_id` FROM `friends` WHERE (sender_id = '%d' AND recipient_id = '%d') OR (sender_id = '%d' AND recipient_id = '%d')", senderId, recipientId, recipientId, senderId)
// 	err := db.QueryRow(selection).Scan(&foundId)

// 	if recipientId == senderId {
// 		return errors.New("cannot send friend request to yourself")
// 	} else if foundId != 0 {
// 		return errors.New("friend request already exists")
// 	} else if err != nil {

// 		insertFriendRequest(fmt.Sprint(senderId), fmt.Sprint(recipientId))

// 		log.Print(fmt.Sprint("recipient : " + recipientEmail))

// 		ping(senderName, recipientEmail, message)
// 		return nil
// 	}

// 	return err
// }
