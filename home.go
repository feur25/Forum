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
	tmpl.ExecuteTemplate(w, "home-page", data)
}
func HandleFriend(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Friend"
	data.Page.Style = "friend"
	getAllFriendRequests(&data.Page.FriendList)
	tmpl.ExecuteTemplate(w, "friends-page", data)
}
