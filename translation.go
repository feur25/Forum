package main

var translation = map[string]map[string]string{
	"accept":  {"FR": "Accepter", "EN": "Accept"},
	"deny":    {"FR": "Refuser", "EN": "Deny"},
	"message": {"FR": "Envoyer un Message", "EN": "Message"},
	"delete":  {"FR": "Supprimer", "EN": "Delete"},
	"send":    {"FR": "Envoyer", "EN": "Send"},
	"welcome": {"FR": "Bienvenue", "EN": "Welcome"},
	"cancel":  {"FR": "Annuler", "EN": "Cancel"},
	"login":   {"FR": "Se Connecter", "EN": "Log in"},
	"logout":  {"FR": "Se Déconnecter", "EN": "Log out"},
	"username":  {"FR": "Pseudonyme", "EN": "Username"},
	"emailaddress":  {"FR": "Adresse Email", "EN": "Username"},
}

func Translate(base string) string {
	val := translation[base]["FR"]
	if val != "" {
		return val
	}
	return "TRANSLATION ERROR"
}
