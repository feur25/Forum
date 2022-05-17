package main

import (
	"log"
	"net/http"
	//"github.com/siongui/godom"
)

// func AcceptOrReject(w http.ResponseWriter, r *http.Request) {
// 	if isButtonPressed(r, "yes") {
// 		acceptFriendRequest("1")
// 	} else if isButtonPressed(r, "no") {
// 		denyFriendRequest("1")
// 	}
// }

func handleHome(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Home"
	data.Page.Style = "home"

	data.Page.Topics = getMostRecentTopics(20)
	data.Page.Friends = getFriendships()
	// AcceptOrReject(w, r)

	log.Print(data.Page.Topics)
	tmpl.ExecuteTemplate(w, "home", data)
}
