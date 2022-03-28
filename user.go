package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func createUserInDatabase(username string, email string, phone string, firstName string, lastName string, address string, date string, password string, admin int) error {
	var datetime = time.Now().UTC().Format("2006-01-02 03:04:05")
	data.time = datetime
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `utilisateur`(`username`, `email`, `phone_number`, `first_name`, `last_name`, `address`, `creation_date`, `password`, `is_admin`) VALUES ('%s','%s','%s','%s','%s','%s', '%s', MD5('%s'), '%d')", username, email, phone, firstName, lastName, address, datetime, password, admin))
	checkError(err)
	defer insert.Close()

	return nil
}

func CheckuserLogin(username, password string) ([]string, error) {
	selectQuery, err := db.Query(fmt.Sprintf("SELECT username, email, password FROM `utilisateur` WHERE username = '%s' OR email = '%s' LIMIT 1", username, username))
	checkError(err)
	defer selectQuery.Close()

	userUsername := ""
	userEmail := ""
	userPassword := ""
	if selectQuery.Next() {
		err := selectQuery.Scan(&userUsername, &userEmail, &userPassword)
		checkError(err)
	}

	if userPassword == password {
		fmt.Println("\ngood Password")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
		})
	} else {
		fmt.Println("\nwrong Password")
		return []string{}, errors.New("wrong Password")
	}

	return nil, nil
}

func checkEmailValidity(email string) bool {
	return strings.Contains(data.email, ".") && strings.Contains(data.email, "@") && strings.Index(data.email, ".") > strings.Index(data.email, "@")
}
func checkPasswordValidity(password string) (bool, error) {
	upper, err := regexp.MatchString(data.password, containsUpperCase)
	checkError(err)
	lower, err := regexp.MatchString(data.password, containsLowerCase)
	checkError(err)
	number, err := regexp.MatchString(data.password, containsNumber)
	checkError(err)
	special, err := regexp.MatchString(data.password, containsSpecial)
	checkError(err)

	if len(data.password) < 10 {
		return false, errors.New("password needs to be 10 characters or more")
	} else if !upper {
		log.Print("Le mot de passe ne contient pas de lettres majuscules.")
		return false, errors.New("password needs to contain upper case character")
	} else if !lower {
		log.Print("Le mot de passe ne contient pas de lettres minuscule.")
		return false, errors.New("password needs to contain lower case character")
	} else if !number {
		log.Print("Le mot de passe ne contient pas de chiffres.")
		return false, errors.New("password needs to contain a number")
	} else if !special {
		log.Print("Le mot de passe ne contient pas de caractères spéciaux.")
		return false, errors.New("password needs to contain special character")
	}

	return true, nil
}
