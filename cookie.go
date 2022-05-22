package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type UserCookie struct {
	id    int
	valid int
}

func InsertCookieUser(validity int) {
	cookie_check, err := db.Query(fmt.Sprintf("INSERT INTO `cookie` (`id`, `cookie_valid`) VALUES('%s', '%d')", data.User.PublicInfo.Id, validity))
	CheckError(err)
	defer cookie_check.Close()
}
func GetCookieUser(id int) (UserCookie, error) {
	cookie := UserCookie{}
	getCookieInformation := fmt.Sprintf("SELECT * FROM cookie WHERE id = '%d' ", id)
	err := db.QueryRow(getCookieInformation).Scan(&cookie.id, &cookie.valid)
	data.User.CookieAccept = UserCookie{}
	data.User.CookieAccept.id = cookie.id
	return cookie, err
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data.User.CookieAccept = UserCookie{}
	cookie, err := r.Cookie("AcceptCookie")
	fmt.Println("cookie:", cookie, "err:", err)
	GetCookieUser(Atoi(data.User.PublicInfo.Id))
	log.Println("cookie : ", data.User.CookieAccept.id)
	if err != nil && data.Login && data.User.CookieAccept.valid != 1 && data.User.CookieAccept.id != Atoi(data.User.PublicInfo.Id) {
		fmt.Println("cookie was not found")
		cookie = &http.Cookie{
			Name:  data.User.PublicInfo.Username,
			Value: data.User.PublicInfo.Id,
			//Domain:   "loaclhost:4448",
			//Path:     "/cookie",
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
		}
		fmt.Println("cookie generate")
		http.SetCookie(w, cookie)
	} else {
		fmt.Println("cookie always active")
	}
	if r.FormValue("AcceptCookie") != "" {
		if data.User.CookieAccept.id != Atoi(data.User.PublicInfo.Id) {
			InsertCookieUser(1)
		}
		Redirect(w, r, "/home")
	}
	if r.FormValue("DenyCookie") != "" {
		if data.User.CookieAccept.id != Atoi(data.User.PublicInfo.Id) {
			InsertCookieUser(0)
		}
		Redirect(w, r, "/home")
	}
	tmpl.ExecuteTemplate(w, "cookie", nil)
}
