package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	PublicInfo  PublicUser
	Password    string
	Email       string
	PhoneNumber string
	FirstName   string
	LastName    string
	Address     string
}
type PublicUser struct {
	Id        string
	Username  string
	ImageLink string
	Admin     int
}

type MessagePrivate struct {
	user_id_friends string
	sender_id       int
	sender          string
	request_friend  bool
	accept          bool
}
type Friend struct {
	user_id string
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func ValidateUserData(user *User) {
	ValidatePublicUserData(&user.PublicInfo)
}
func ValidatePublicUserData(user *PublicUser) {
	if user.ImageLink == "" {
		user.ImageLink = "/images/defaultpfp.jpg"
	}
}

func createUserInDB(username string, email string, phone string, firstName string, lastName string, address string, password string, admin int) error {
	var datetime = time.Now().UTC().Format("2006-01-02 03:04:05")

	var alreadyExists bool
	err := db.QueryRow(fmt.Sprintf("SELECT IF(COUNT(*),'true','false') FROM users WHERE email = '%s' ", email)).Scan(&alreadyExists)
	checkError(err)
	if alreadyExists {
		return errors.New("user with email already exists")
	}

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `users` (`username`, `email`, `phone_number`, `first_name`, `last_name`, `address`, `creation_date`, `password`, `is_admin`) VALUES ('%s','%s','%s','%s','%s','%s', '%s', '%s', '%d')", username, email, phone, firstName, lastName, address, datetime, MD5(password), admin))
	checkError(err)
	defer insert.Close()

	return nil
}
func getUserInDB(user string) (User, error) {
	var id, username, imgLink, pass, email, phone, lastN, firstN, address, date string
	var admin int
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE email = '%s' OR username = '%s' ", user, user)).Scan(&id, &username, &imgLink, &pass, &email, &admin, &phone, &firstN, &lastN, &address, &date)

	auth := User{PublicUser{id, username, imgLink, admin}, pass, email, phone, firstN, lastN, address}
	ValidateUserData(&auth)
	return auth, err
}
func getUserWithIdInDB(userId string) (User, error) {
	var id, username, imgLink, pass, email, phone, lastN, firstN, address, date string
	var admin int
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE user_id = '%s' ", userId)).Scan(&id, &username, &imgLink, &pass, &email, &admin, &phone, &firstN, &lastN, &address, &date)

	auth := User{PublicUser{id, username, imgLink, admin}, pass, email, phone, firstN, lastN, address}
	ValidateUserData(&auth)
	return auth, err
}

func updateUserInDB(w http.ResponseWriter, r *http.Request, password string) error {
	if data.User.PublicInfo.Id != "" {
		update, err := db.Query(fmt.Sprintf("UPDATE `users` SET `password`='%s' WHERE `username`='%s'", password, data.User.PublicInfo.Username))
		checkError(err)
		defer update.Close()
	} else {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
	}
	return nil
}
func deleteUserInDB(username string, password string) error {
	delete, err := db.Query(fmt.Sprintf("DELETE FROM `users` WHERE `username`='%s' AND `password`='%s'", username, password))
	checkError(err)
	defer delete.Close()

	return nil
}
func addFriendInDB(id string) (MessagePrivate, error) {
	var id_other, sender string
	var request, accept bool
	var id_sender int
	errors := db.QueryRow(fmt.Sprintf("SELECT * FROM private_sender WHERE user_id = '%s'", id)).Scan(&id_other, &id_sender, &sender, &request, &accept)
	checkError(errors)
	Identifiant := MessagePrivate{id_other, id_sender, sender, request, accept}
	addfriend, err := db.Query(fmt.Sprintf("INSERT INTO `friends` (`user_id`, `user_id_1`) VALUES ('%s','%s')", id, id_other))
	checkError(err)
	defer addfriend.Close()
	return Identifiant, nil
}

func privateSenderInDB(id string, sender string, request bool, accept bool) error {
	if !request {
		log.Printf(data.Message.sender)

	} else if request {
		data.Message.sender += "\n" + data.User.PublicInfo.Username + " Souhaite devenir votre amigo !!"
	}
	message, err := db.Query(fmt.Sprintf("INSERT INTO `private_sender` (`user_id`, `sender`, `request_friend`, `accept`) VALUES ('%s','%s', '%t', '%t' )", id, sender, request, accept))
	checkError(err)
	defer message.Close()
	return nil
}
func checkFriendInDB(id string) (Friend, error) {
	var id1, id2 string
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM friends WHERE user_id_1 = '%s'", id)).Scan(&id1, &id2)
	checkError(err)
	check := Friend{id2}
	return check, nil
}

// Admin : AdminAdmin1234567890/
