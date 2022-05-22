package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	PublicInfo   PublicUser
	CookieAccept UserCookie
	PingName     Ping
	Password     string
	Email        string
	PhoneNumber  string
	FirstName    string
	LastName     string
	Address      string
	CreationTime string
}
type PublicUser struct {
	Id        string
	Username  string
	ImageLink string
	Admin     int
}
type Ping struct {
	name string
}

type MessagePrivate struct {
	FriendId         int
	SenderId         int
	SenderContent    string
	date             string
	RecipientContent string
	RecipientId      int
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

func createUser(username string, email string, phone string, firstName string, lastName string, address string, password string, admin int) error {

	var alreadyExists bool
	err := db.QueryRow(fmt.Sprintf("SELECT IF(COUNT(*),'true','false') FROM users WHERE email = '%s' ", email)).Scan(&alreadyExists)
	CheckError(err)
	if alreadyExists {
		return errors.New("user with email already exists")
	}

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `users` (`username`, `email`, `phone_number`, `first_name`, `last_name`, `address`, `creation_date`, `password`, `is_admin`) VALUES ('%s','%s','%s','%s','%s','%s', '%s', '%s', '%d')", username, email, phone, firstName, lastName, address, datetime, MD5(password), admin))
	CheckError(err)
	defer insert.Close()

	return nil
}

func getUserWithUsername(userId string) (User, error) {
	log.Print(userId)
	user := User{}
	selection := fmt.Sprintf("SELECT * FROM users WHERE email = '%s' OR username = '%s' ", userId, userId)
	err := db.QueryRow(selection).Scan(&user.PublicInfo.Id, &user.PublicInfo.ImageLink, &user.PublicInfo.Username, &user.Password, &user.Email, &user.PublicInfo.Admin, &user.PhoneNumber, &user.FirstName, &user.LastName, &user.Address, &user.CreationTime)
	ValidateUserData(&user)
	return user, err
}
func getUserWithId(userId string) (User, error) {
	user := User{}
	selection := fmt.Sprintf("SELECT * FROM users WHERE user_id = '%s' ", userId)
	err := db.QueryRow(selection).Scan(&user.PublicInfo.Id, &user.PublicInfo.ImageLink, &user.PublicInfo.Username, &user.Password, &user.Email, &user.PublicInfo.Admin, &user.PhoneNumber, &user.FirstName, &user.LastName, &user.Address, &user.CreationTime)
	ValidateUserData(&user)
	return user, err
}

func updateUser(w http.ResponseWriter, r *http.Request, password string) error {
	if data.User.PublicInfo.Id != "" {
		update, err := db.Query(fmt.Sprintf("UPDATE `users` SET `password`='%s' WHERE `username`='%s'", password, data.User.PublicInfo.Username))
		CheckError(err)
		defer update.Close()
	} else {
		Redirect(w, r, "/login")
	}
	return nil
}
func deleteUser(username string, password string) error {
	delete, err := db.Query(fmt.Sprintf("DELETE FROM `users` WHERE `username`='%s' AND `password`='%s'", username, password))
	CheckError(err)
	defer delete.Close()

	return nil
}

func idToUsername(id string) string {
	selection := fmt.Sprintf("SELECT username FROM users WHERE user_id = '%s'", id)
	username := ""
	db.QueryRow(selection).Scan(&username)
	return username
}
func CheckIfUserExist(name string) bool {
	selection, err := db.Query(fmt.Sprintf("SELECT username FROM users WHERE name = '%s'", name))
	if err != nil {
		log.Println("La personne n'existe pas !")
	}
	defer selection.Close()
	data.User.PingName.name = name
	return true
}

// Admin : AdminPassword1234/
