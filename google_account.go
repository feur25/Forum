package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserGoogleInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}
type UserInfo struct {
	Google     UserGoogleInfo ` json:"UserInfo"`
	UserGoogle GoogleUser
}

func FetchUserInfo(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:4448/auth/google/callback",
	ClientID:     "782907329281-43c0m92f188c6enhk15e1jfvg6gcn4fh.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleUser = GoogleUser{}
var responses = UserGoogleInfo{}
var test UserGoogleInfo = UserGoogleInfo{}

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w)
	log.Println("Nom : ")
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	createUser(googleUser.Name, googleUser.Email, "", googleUser.GivenName, googleUser.FamilyName, "", "Zeusdu96!F", 0)
	log.Println("let's go 1!!")
	//tmpl.ExecuteTemplate(w, "googlev2", nil)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	data, _ := getUserDataFromGoogle(r.FormValue("code"))
	email := r.FormValue("email")
	password := r.FormValue("password")
	pseudo := r.FormValue("pseudo")
	name := r.FormValue("name")
	lastname := r.FormValue("lastname")
	if r.FormValue("google") != "" {
		createUserGoogle(responses.Picture, pseudo, email, "", name, lastname, "", password, 0)
		log.Println(pseudo, email, password, name, lastname)
		log.Println("Vous avez terminer : ", responses.Picture)
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
	}
	log.Println(w, "UserInfo: %s\n", data)
	//var googleUser GoogleUser = GoogleUser{}
	//createUser(googleUser.Name, responses.Email, "", googleUser.GivenName, googleUser.FamilyName, "", "Zeusdu96!F", 0)
	log.Println("let's go!!")
	tmpl.ExecuteTemplate(w, "google", nil)
}
func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	log.Println("Nom6 : ", GoogleUser{})
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	j := json.Unmarshal(contents, &responses)
	log.Println("plop : ", j)
	log.Println("plop : ", test.Email)
	return contents, nil
}
