package main

import (
	"net/http"
	//"github.com/siongui/godom"
)

// func AcceptOrReject(w http.ResponseWriter, r *http.Request) {
// 	if IsButtonPressed(r, "yes") {
// 		acceptFriendRequest("1")
// 	} else if IsButtonPressed(r, "no") {
// 		denyFriendRequest("1")
// 	}
// }

func HandleHome(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Home"
	data.Page.Style = "home"

	data.Page.TopicList.Topics = getMostRecentTopics(20)
	getAllFriendRequests(&data.Page.FriendList)
	// AcceptOrReject(w, r)

	// log.Print(data.Page.Topics)
	tmpl.ExecuteTemplate(w, "home-page", data)
}
