package main

import (
	"fmt"
	"log"
	"net/http"
)

// func handleGetCode(w http.ResponseWriter, r *http.Request) {
// 	data.Page.Title = "code"
// 	data.Page.Style = "code"
// 	code := r.FormValue("code")
// 	if r.FormValue("code_drop") == "getcode" {
// 		data.UpdateConfirmationCode = send_update_email(data.User.email)
// 		log.Print("\ncode : " + data.UpdateConfirmationCode + "\n")
// 	}
// 	if code == data.UpdateConfirmationCode && len(data.UpdateConfirmationCode) >= 5 {
// 		data.Authorized = true
// 		if data.Authorized {
// 			http.Redirect(w, r, "http://"+Host+":"+Port+"/update", http.StatusMovedPermanently)
// 		}
// 	}
// 	tmpl.ExecuteTemplate(w, "code", data)
// }

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if !data.Login {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
		return
	}

	data.Page.Title = "Update"
	data.Page.Style = "update"

	if data.UpdateConfirmationCode == "" {
		data.UpdateConfirmationCode = send_update_email(data.User.Email)
	}

	if r.FormValue("envoyer") == "Envoyer" && r.FormValue("code") == data.UpdateConfirmationCode {
		password := r.FormValue("password")
		passwordcheck := r.FormValue("checkpassword")

		log.Print(passwordcheck + " : " + password)

		if password != passwordcheck {
			var bad_password = `<p class="error-message">Les deux mots de passe ne correspondent pas</p>`
			fmt.Fprint(w, bad_password)
		} else if MD5(password) == data.User.Password {
			var password_active = `<p class="error-message">Le nouveau mot de passe ressemble a votre ancien mot de passe</p>`
			fmt.Fprint(w, password_active)
		} else {
			log.Print("Le mot de passe a bien changé")
			updateUserInDB(w, r, MD5(password))
			data.User.Password = password
			http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)

			data.UpdateConfirmationCode = ""
		}
	}

	tmpl.ExecuteTemplate(w, "update", data)

	// } else {
	// 	var bad_code_message = `<p style="color: red;">Vous devez d'abord récuperer votre code de vérification, avant de vouloir réinitialiser votre mot de passe !`
	// 	fmt.Fprint(w, bad_code_message)
	// }
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if !data.Login {
		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
		return
	}

	data.Page.Title = "Delete"
	data.Page.Style = "delete"

	if data.DeleteConfirmationCode == "" {
		data.DeleteConfirmationCode = send_delete_email(data.User.Email)
	}

	if r.FormValue("envoyer") == "Envoyer" && r.FormValue("code") == data.DeleteConfirmationCode {
		log.Print(data.User.Password)

		http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
		deleteUserInDB(data.User.PublicInfo.Username, data.User.Password)

		data.DeleteConfirmationCode = ""
	}
	tmpl.ExecuteTemplate(w, "delete", data)
}

/*func profile(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Profile"
	data.Page.Style = "profile"
	if data.User.PublicInfo.user_id == data.Message.user_id_friends {
		tmpl.ExecuteTemplate(w, "profile", data)
	} else {
		tmpl.ExecuteTemplate(w, "profile2", data)
	}
}*/
func CheckSender() {
	switch {
	case data.friend_info.user_id == "" && !data.button:
		privateSenderInDB(9, data.Message.sender, 0, 0)
	case data.friend_info.user_id == "" && data.button:
		privateSenderInDB(9, data.Message.sender, 1, 0)
	case data.friend_info.user_id != "":
		privateSenderInDB(9, data.Message.sender, 0, 1)
	}
}
func HandleSenderMessage(w http.ResponseWriter, r *http.Request) {
	data.Page.Title = "Call"
	data.Page.Style = "call"
	if data.Login {
		//checkFriendInDB(data.User.PublicInfo.Id)
		text := r.FormValue("message")
		if r.FormValue("envoyer") == "Envoyer" {
			data.Message.sender = text
			CheckSender()
			privateSenderInDB(9, data.Message.sender, 0, 0)
		}
	}
	tmpl.ExecuteTemplate(w, "call", data)
}
