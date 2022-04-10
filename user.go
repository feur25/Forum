package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func MD5(raw string) string {
	encrypted := md5.Sum([]byte(raw))
	return hex.EncodeToString(encrypted[:])
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
func getUserInDB(user string) (Auth, error) {
	var id, username, pass, email, phone, lastN, firstN, address, date string
	var admin bool
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE email = '%s' OR username = '%s' ", user, user)).Scan(&id, &username, &pass, &email, &admin, &phone, &firstN, &lastN, &address, &date)

	auth := Auth{id, username, pass, email, phone, firstN, lastN, address}
	return auth, err
}
func getUserWithIdInDB(userId string) (Auth, error) {
	var id, username, pass, email, phone, lastN, firstN, address, date string
	var admin bool
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE user_id = '%s' ", userId)).Scan(&id, &username, &pass, &email, &admin, &phone, &firstN, &lastN, &address, &date)

	auth := Auth{id, username, pass, email, phone, firstN, lastN, address}
	return auth, err
}

func updateUserInDB(w http.ResponseWriter, r *http.Request, password string) error {
	if data.Auth.user_id != "" {
		update, err := db.Query(fmt.Sprintf("UPDATE `users` SET `password`='%s'", password))
		checkError(err)
		defer update.Close()
	} else {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
	}
	return nil
}
func deleteUserInDB(username string, password string) error {
	delete, err := db.Query(fmt.Sprintf("DELETE FROM `users` WHERE `username`='%s', AND `password`='%s' ", username, password))
	checkError(err)
	defer delete.Close()

	return nil
}

// Admin : AdminAdmin1234567890/
