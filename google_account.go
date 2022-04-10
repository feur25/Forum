package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// "oauth-example/controller"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func googleLogin(w http.ResponseWriter, r *http.Request) {
	/*googleConfig := SetupConfig()
	url := googleConfig.AuthCodeURL("randomstate")*/
	http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently)

}

func googleCallback(r http.ResponseWriter, w *http.Request) {
	/*state := w.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Fprintln(r, "states dont match")
		return

	}*/
	code := w.URL.Query()["code"][0]
	googleConfig := setupConfig()
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintln(r, "merde")

	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	checkError(err)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	checkError(err)

	var googleData GoogleUser = GoogleUser{}
	err = json.Unmarshal(contents, &googleData)
	checkError(err)

	fmt.Fprintln(r, string(contents))
}

func setupConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     "782907329281-5g05saomlg8p20r9214k6cg3e2sirs41.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8888/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
