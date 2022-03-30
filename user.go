package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func MD5(raw string) string {
	encrypted := md5.Sum([]byte(raw))
	return hex.EncodeToString(encrypted[:])
}

func createUserInDatabase(username string, email string, phone string, firstName string, lastName string, address string, password string, admin int) error {
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

func checkUserLogin(w http.ResponseWriter, r *http.Request, username, password string) ([]string, error) {
	selectQuery, err := db.Query(fmt.Sprintf("SELECT username, email, password, user_id FROM `users` WHERE username = '%s' OR email = '%s' LIMIT 1", username, username))
	checkErrorLogout(err)
	defer selectQuery.Close()

	userUsername := ""
	userEmail := ""
	userPassword := ""
	userId := ""
	if selectQuery.Next() {
		err := selectQuery.Scan(&userUsername, &userEmail, &userPassword, &userId)
		checkErrorLogout(err)
	}
	log.Println(userPassword + MD5(password))
	if userPassword == MD5(password) {
		fmt.Println("\ngood Password")
		data.user_id = userId
		log.Println(data.user_id)
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
	} else {
		fmt.Println("\nwrong Password")
		return []string{}, errors.New("wrong Password")
	}

	return nil, nil
}

func checkUsernameValidity(username string) (bool, error) {
	if len(username) < 4 {
		return false, errors.New("username must be at least 4 characters long")
	}
	return true, nil
}
func checkAdressValidity(adress string) (bool, error) {
	if len(adress) < 4 {
		return false, errors.New("address must not be empty")
	}
	return true, nil
}
func checkEmailValidity(email string) (bool, error) {
	if !strings.Contains(email, ".") || !strings.Contains(email, "@") || strings.Index(email, ".") <= strings.Index(email, "@") || len(email) < 3 {
		return false, errors.New("email must follow the standard email structure : \"example@website.com\"")
	}
	return true, nil
}
func checkPasswordValidity(w http.ResponseWriter, r *http.Request, password string) (bool, error) {

	fmt.Println(password)
	err := checkPasswordErrors(password)

	return err == nil, err
}

func checkPasswordErrors(password string) error {
	switch {
	case len(password) < 10:
		return errors.New("the password has to be 10 characters or more")
	case strings.ToLower(password) == password:
		return errors.New("the password has to contain upper case character")
	case strings.ToUpper(password) == password:
		return errors.New("the password has to contain lower case character")
	case !strings.ContainsAny(password, "1234567890"):
		return errors.New("the password has to contain a number")
	case !strings.ContainsAny(password, `-+_!@#$%^&*.,?/\`):
		return errors.New("the password has to contain a special character")
	case strings.Contains(password, ` `):
		return errors.New("the password has to not contain any space")
	default:
		return nil
	}
}

func updateUserInDatabase(password string) error {
	update, err := db.Query(fmt.Sprintf("UPDATE `users` SET `password`='%s'", password))
	checkError(err)
	defer update.Close()

	return nil
}
func deleteUserInDatabase(username string, password string) error {
	delete, err := db.Query(fmt.Sprintf("DELETE FROM `users` WHERE `username`='%s', AND `password`='%s' ", username, password))
	checkError(err)
	defer delete.Close()

	return nil
}
