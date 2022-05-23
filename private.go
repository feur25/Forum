package main

import (
	"fmt"
	"log"
	"time"
)

func PrivateMessage(GroupId int, SenderId int, RecipientId int, SenderContent string, RecipientContent string) error {
	var DateTimeSenderMessage = time.Now().UTC().Format("2006-01-02 03:04:05")
	AddMessagePrivate, err := db.Query(fmt.Sprintf("INSERT INTO `message` (`friend_id`, `sender_id`, `recipient_id`, `sender_content`, `recipient_content`, `date`) VALUES ('%d','%d','%d','%s','%s','%s')", GroupId, SenderId, RecipientId, SenderContent, RecipientContent, DateTimeSenderMessage))
	CheckError(err)
	defer AddMessagePrivate.Close()

	return nil
}

/*func GetSenderId(sender_id int) (MessagePrivate, error) {
	sender := MessagePrivate{}
	info := fmt.Sprintf("SELECT * FROM message WHERE sender_id = '%d'", sender_id)
	err := db.QueryRow(info).Scan(&sender.FriendId, &sender.SenderId, &sender.RecipientId, &sender.SenderContent, &sender.RecipientContent, &sender.date)
	return sender, err
}*/
func GetMesssage(id int) (MessagePrivate, error) {
	recipient := MessagePrivate{}
	info := fmt.Sprintf("SELECT * FROM message WHERE recipient_id = '%d' OR sender_id = '%d'", id, id)
	err := db.QueryRow(info).Scan(&recipient.FriendId, &recipient.SenderId, &recipient.RecipientId, &recipient.SenderContent, &recipient.RecipientContent, &recipient.date)
	return recipient, err
}
func DisplayMessage() (MessagePrivate, error) {
	id := Atoi(data.User.PublicInfo.Id)
	sender, err := GetMesssage(id)
	if id == data.Message.SenderId {

	} else {
		log.Println("You not owner in this sender")
	}
	return sender, err
}
