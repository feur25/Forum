package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func HandleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Call"
	data.Page.Style = "call"

	//checkFriendInDB(data.User.PublicInfo.Id)
	text := r.FormValue("message")
	// caca := godom.Document.QuerySelector("name=()")
	if isButtonPressed(r, "envoyer") {
		if !data.Login {
			http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
		} else {
			senderId, _ := strconv.Atoi(data.User.PublicInfo.Id)
			recipientUsername := r.FormValue("id")
			err := createFriendRequest(senderId, recipientUsername, text)
			checkError(err)
		}
	}

	tmpl = tmpl.Funcs(template.FuncMap{
		"AcceptFriendRequest": AcceptFriendRequest,
		"DenyFriendRequest":   DenyFriendRequest,
	})

	tmpl.ExecuteTemplate(w, "call", data)
}

func AcceptFriendRequest(id string) {
	log.Print("Accepted " + id)
	// edit := fmt.Sprintf("UPDATE friends SET pending = '0' WHERE friends.request_id = %s", id)
	// db.Query(edit)
}

func DenyFriendRequest(id string) {
	log.Print("Denied " + id)
	// edit := fmt.Sprintf("DELETE FROM friends WHERE request_id = %s", id)
	// db.Query(edit)
}

func getFriendships() FriendList {
	selection := fmt.Sprintf("SELECT request_id, sender_id, recipient_id, pending FROM friends WHERE sender_id = '%s' OR recipient_id = '%s' ", data.User.PublicInfo.Id, data.User.PublicInfo.Id)
	query, _ := db.Query(selection)
	defer query.Close()

	result := FriendList{}
	for query.Next() {
		var friend Friend
		query.Scan(&friend.FriendId, &friend.SenderId, &friend.RecipientId, &friend.Pending)

		if friend.SenderId == data.User.PublicInfo.Id {
			friend.FriendName = idToUsername(friend.RecipientId)
		} else {
			friend.FriendName = idToUsername(friend.SenderId)
		}
		if friend.Pending == 1 {
			if friend.SenderId == data.User.PublicInfo.Id {
				result.SentRequests = append(result.SentRequests, friend)
			} else {
				result.ReceivedRequests = append(result.ReceivedRequests, friend)
			}
		} else {
			result.AcceptedRequests = append(result.AcceptedRequests, friend)
		}
	}
	return result
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

		insertion := fmt.Sprintf("INSERT INTO `friends` (`sender_id`, `recipient_id`) VALUES ('%d','%d' )", senderId, recipientId)
		insertQuery, err := db.Query(insertion)
		checkError(err)
		defer insertQuery.Close()

		log.Print(fmt.Sprint("recipient : " + recipientEmail))

		ping(senderName, recipientEmail, message)
		return nil
	}

	return err
}
