package main

import (
	"errors"
	"fmt"

	"net/http"
	"strings"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {

	data.Page.Title = "Register"
	data.Page.Style = "register"

	username := r.FormValue("username")
	password := r.FormValue("password")
	address := r.FormValue("address")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	submit := r.FormValue("submit")
	if submit != "" {

		usernameValidity, err := checkUsernameValidity(username)
		checkErrorLogout(err)
		addressValidity, err := checkAdressValidity(address)
		checkErrorLogout(err)
		emailValidity, err := checkEmailValidity(email)
		checkErrorLogout(err)
		passValidity, err := checkPasswordValidity(w, r, password)
		checkErrorLogout(err)

		if usernameValidity && addressValidity && len(phone) >= 10 && emailValidity && passValidity && len(first_name) >= 4 && len(last_name) >= 3 {

			err := createUserInDB(username, email, phone, first_name, last_name, address, password, normal_user)
			if !checkErrorLogout(err) {
				registerSuccess(w, r, email)
			} else {
				registerFail(w, r)
			}
		}

	}

	submit = ""
	tmpl.ExecuteTemplate(w, "register", data)
}

func registerSuccess(w http.ResponseWriter, r *http.Request, email string) {
	http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)
}

func registerFail(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+Host+":"+Port+"/register", http.StatusMovedPermanently)
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
