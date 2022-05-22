package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type FriendList struct {
	SelectedFriend  Friend
	AcceptedFriends []Friend
	SentFriends     []Friend
	ReceivedFriends []Friend
}

type Friend struct {
	RequestId   string
	SenderId    string
	RecipientId string
	Pending     int
	FriendUser  PublicUser
}

func HandleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Call"
	data.Page.Style = "call"

	//checkFriend(data.User.PublicInfo.Id)
	text := r.FormValue("message")
	// caca := godom.Document.QuerySelector("name=()")
	if IsButtonPressed(r, "envoyer") {
		if !data.Login {
			http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
		} else {
			senderId, _ := strconv.Atoi(data.User.PublicInfo.Id)
			recipientUsername := r.FormValue("id")
			err := createFriendRequest(senderId, recipientUsername, text)
			CheckError(err)
		}
	}

	// tmpl = tmpl.Funcs(template.FuncMap{
	// 	"AcceptFriendRequest": AcceptFriendRequest,
	// 	"DenyFriendRequest":   DenyFriendRequest,
	// })

	tmpl.ExecuteTemplate(w, "call", data)
}
func checkFriendAuthorization(requestId, thisId string) bool {
	var authorized bool
	selection := fmt.Sprintf("SELECT IF(COUNT(*),'true','false') FROM friends WHERE request_id = '%s' AND (sender_id = '%s' OR recipient_id = '%s') ", requestId, thisId, thisId)
	db.QueryRow(selection).Scan(&authorized)
	if !authorized {
		log.Print("Unauthorized access to friend operation by : " + thisId)
	}
	return authorized
}

func checkSenderFriendAuthorization(requestId, thisId string) bool {
	var authorized bool
	selection := fmt.Sprintf("SELECT IF(COUNT(*),'true','false') FROM friends WHERE request_id = '%s' AND sender_id = '%s' ", requestId, thisId)
	db.QueryRow(selection).Scan(&authorized)
	if !authorized {
		log.Print("Unauthorized access to friend operation by : " + thisId)
	}
	return authorized
}

func checkRecipientFriendAuthorization(requestId, thisId string) bool {
	var authorized bool
	selection := fmt.Sprintf("SELECT IF(COUNT(*),'true','false') FROM friends WHERE request_id = '%s' AND recipient_id = '%s' ", requestId, thisId)
	db.QueryRow(selection).Scan(&authorized)
	if !authorized {
		log.Print("Unauthorized access to friend operation by : " + thisId)
	}
	return authorized
}

func HandleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	defer Redirect(w, r, "/home")

	requestId, err := GetUrlParam(r, "id")
	if CheckError(err) {
		log.Print(err)
	}

	if checkRecipientFriendAuthorization(requestId, data.User.PublicInfo.Id) {
		edit := fmt.Sprintf("UPDATE friends SET pending = '0' WHERE friends.request_id = %s", requestId)
		db.Query(edit)
	}

	// log.Print(requestId)
}

func HandleDenyFriendRequest(w http.ResponseWriter, r *http.Request) {
	log.Print("caca")
	defer Redirect(w, r, "/home")

	requestId, err := GetUrlParam(r, "id")
	if CheckError(err) {
		log.Print(err)
	}
	log.Print("salope de pute")
	if checkFriendAuthorization(requestId, data.User.PublicInfo.Id) {
		edit := fmt.Sprintf("DELETE FROM friends WHERE request_id = %s AND pending = '1'", requestId)
		db.Query(edit)
	}
}

func HandleDeleteFriend(w http.ResponseWriter, r *http.Request) {
	defer Redirect(w, r, "/home")

	requestId, err := GetUrlParam(r, "id")
	if CheckError(err) {
		log.Print(err)
	}

	if checkFriendAuthorization(requestId, data.User.PublicInfo.Id) {
		edit := fmt.Sprintf("DELETE FROM friends WHERE request_id = %s AND pending = '0'", requestId)
		db.Query(edit)
	}
}

func HandleMessageFriend(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Topic"
	data.Page.Style = "topic"

	topicId, err := GetUrlParam(r, "id")
	if !checkFriendAuthorization(topicId, data.User.PublicInfo.Id) {
		Redirect(w, r, "/home")
		return
	}

	if CheckError(err) {
		data.Page.FriendList.SelectedFriend = Friend{}
	} else {
		data.Page.FriendList.SelectedFriend = getFriendRequest(topicId)
	}

	tmpl.ExecuteTemplate(w, "message", data)
}

func getFriendships(list *FriendList) {
	selection := fmt.Sprintf("SELECT friends.*, sender.user_id, sender.username, sender.image_link, sender.is_admin, recipient.user_id, recipient.username, recipient.image_link, recipient.is_admin FROM friends LEFT JOIN users as sender ON sender.user_id = friends.sender_id LEFT JOIN users as recipient ON recipient.user_id = friends.recipient_id WHERE sender_id = '%s' OR recipient_id = '%s'", data.User.PublicInfo.Id, data.User.PublicInfo.Id)
	query, _ := db.Query(selection)
	defer query.Close()

	accepted := []Friend{}
	sent := []Friend{}
	received := []Friend{}
	for query.Next() {
		friend := Friend{}
		sender := PublicUser{}
		recipient := PublicUser{}
		query.Scan(&friend.RequestId, &friend.SenderId, &friend.RecipientId, &friend.Pending, &sender.Id, &sender.Username, &sender.ImageLink, &sender.Admin, &recipient.Id, &recipient.Username, &recipient.ImageLink, &recipient.Admin)
		ValidatePublicUserData(&sender)
		ValidatePublicUserData(&recipient)
		if data.User.PublicInfo.Id == sender.Id {
			friend.FriendUser = recipient
		} else {
			friend.FriendUser = sender
		}

		if friend.Pending == 0 {
			accepted = append(accepted, friend)
		} else if data.User.PublicInfo.Id == sender.Id {
			sent = append(sent, friend)
		} else {
			received = append(received, friend)
		}
	}
	list.AcceptedFriends = accepted
	list.SentFriends = sent
	list.ReceivedFriends = received
}

func getFriendRequest(requestId string) Friend {
	friend := Friend{}
	insertion := fmt.Sprintf("SELECT * FROM `friends` `sender_id` = '%s'", requestId)
	db.QueryRow(insertion).Scan(&friend.RequestId, &friend.RecipientId, &friend.RecipientId, &friend.Pending)
	if friend.SenderId == data.User.PublicInfo.Id {
		friend.FriendUser = data.User.PublicInfo
	} else {
		user, _ := getUserWithId(friend.SenderId)
		friend.FriendUser = user.PublicInfo
	}
	return friend
}

func insertFriendRequest(senderId, recipientId string) Friend {
	insertion := fmt.Sprintf("INSERT INTO `friends` (`sender_id`, `recipient_id`) VALUES ('%s','%s' )", senderId, recipientId)
	db.QueryRow(insertion)

	friend := Friend{}
	selection := fmt.Sprintf("SELECT * FROM `friends` `sender_id` = '%s' AND `recipient_id` = '%s'", senderId, recipientId)
	db.QueryRow(selection).Scan(&friend.RequestId, &friend.SenderId, &friend.RecipientId, &friend.Pending)
	return friend
}

func createFriendRequest(senderId int, recipientName string, message string) error {
	var recipientId int
	var recipientEmail, senderName string
	selection := fmt.Sprintf("SELECT recipient.user_id, recipient.email, sender.username FROM users as recipient, users as sender WHERE sender.user_id = %d AND recipient.username = '%s'", senderId, recipientName)
	db.QueryRow(selection).Scan(&recipientId, &recipientEmail, &senderName)

	foundId := 0
	selection = fmt.Sprintf("SELECT `request_id` FROM `friends` WHERE (sender_id = '%d' AND recipient_id = '%d') OR (sender_id = '%d' AND recipient_id = '%d')", senderId, recipientId, recipientId, senderId)
	err := db.QueryRow(selection).Scan(&foundId)

	if recipientId == senderId {
		return errors.New("cannot send friend request to yourself")
	} else if foundId != 0 {
		return errors.New("friend request already exists")
	} else if err != nil {

		insertFriendRequest(fmt.Sprint(senderId), fmt.Sprint(recipientId))

		log.Print(fmt.Sprint("recipient : " + recipientEmail))

		ping(senderName, recipientEmail, message)
		return nil
	}

	return err
}
