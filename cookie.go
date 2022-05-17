package main

import (
	"fmt"
	"net/http"
	"time"
)

func isExpired() bool {
	var s = session{}
	return s.expiry.Before(time.Now())
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("accept")
	fmt.Println("cookie:", cookie, "err:", err)
	if err != nil && data.Login {
		fmt.Println("cookie was not found")
		cookie = &http.Cookie{
			Name:     data.User.PublicInfo.Username,
			Value:    data.User.PublicInfo.Id,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
	tmpl.ExecuteTemplate(w, "cookie", nil)
}
