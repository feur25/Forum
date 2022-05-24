package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type FriendList struct {
	SelectedFriend         Friend
	SelectedFriendMessages []FriendMessage
	AcceptedFriends        []Friend
	SentFriends            []Friend
	ReceivedFriends        []Friend
}

type Friend struct {
	RequestId   string
	SenderId    string
	RecipientId string
	Pending     int
	FriendUser  PublicUser
}

func HandleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
	recipientUsername, err := GetUrlParam(r, "name")

	if !CheckError(err) {
		if !data.Login {
			TemporaryRedirect(w, r, "/login")
			return
		}
		log.Print(data.User.PublicInfo.Id + " Sending friend request to : " + recipientUsername)
		err := sendFriendRequest(data.User.PublicInfo.Id, recipientUsername, "Salope")
		CheckError(err)
		log.Print("Hello")
		TemporaryRedirect(w, r, "/friends")
	}
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
	defer TemporaryRedirect(w, r, "/friends")

	requestId, err := GetUrlParam(r, "id")
	if CheckError(err) {
		log.Print(err)
	}

	if checkRecipientFriendAuthorization(requestId, data.User.PublicInfo.Id) {
		edit := fmt.Sprintf("UPDATE friends SET pending = '0' WHERE friends.request_id = %s", requestId)
		db.Query(edit)
	}
}

func HandleDenyFriendRequest(w http.ResponseWriter, r *http.Request) {
	defer TemporaryRedirect(w, r, "/friends")

	requestId, err := GetUrlParam(r, "id")
	if CheckError(err) {
		log.Print(err)
	}
	if checkFriendAuthorization(requestId, data.User.PublicInfo.Id) {
		edit := fmt.Sprintf("DELETE FROM friends WHERE request_id = %s AND pending = '1'", requestId)
		db.Query(edit)
	}
}

func HandleDeleteFriend(w http.ResponseWriter, r *http.Request) {
	defer TemporaryRedirect(w, r, "/friends")

	requestId, err := GetUrlParam(r, "id")
	if CheckError(err) {
		log.Print(err)
	}

	if checkFriendAuthorization(requestId, data.User.PublicInfo.Id) {
		edit := fmt.Sprintf("DELETE FROM friends WHERE request_id = %s AND pending = '0'", requestId)
		db.Query(edit)
	}
}

func getAllFriendRequests(list *FriendList) {
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
	selection := fmt.Sprintf("SELECT * FROM `friends` WHERE `request_id` = '%s' ", requestId)
	db.QueryRow(selection).Scan(&friend.RequestId, &friend.SenderId, &friend.RecipientId, &friend.Pending)
	if friend.SenderId == data.User.PublicInfo.Id {
		user, _ := getUserWithId(friend.RecipientId)
		friend.FriendUser = user.PublicInfo
	} else {
		user, _ := getUserWithId(friend.SenderId)
		friend.FriendUser = user.PublicInfo
	}
	log.Print(friend)
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

func sendFriendRequest(senderId, recipientName string, message string) error {
	var recipientId, recipientEmail, senderName string
	selection := fmt.Sprintf("SELECT recipient.user_id, recipient.email, sender.username FROM users as recipient, users as sender WHERE sender.user_id = %s AND recipient.username = '%s'", senderId, recipientName)
	db.QueryRow(selection).Scan(&recipientId, &recipientEmail, &senderName)

	foundId := 0
	selection = fmt.Sprintf("SELECT `request_id` FROM `friends` WHERE (sender_id = '%s' AND recipient_id = '%s') OR (sender_id = '%s' AND recipient_id = '%s')", senderId, recipientId, recipientId, senderId)
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
