package main

var translation = map[string]map[string]string{
	"accept":       {"EN": "Accept", "FR": "Accepter"},
	"deny":         {"EN": "Deny", "FR": "Refuser"},
	"sendmessage":  {"EN": "Message", "FR": "Envoyer un Message"},
	"delete":       {"EN": "Delete", "FR": "Supprimer"},
	"send":         {"EN": "Send", "FR": "Envoyer"},
	"welcome":      {"EN": "Welcome", "FR": "Bienvenue"},
	"cancel":       {"EN": "Cancel", "FR": "Annuler"},
	"login":        {"EN": "Log in", "FR": "Se Connecter"},
	"logout":       {"EN": "Log out", "FR": "Se Déconnecter"},
	"register":     {"EN": "Register", "FR": "S'inscrire"},
	"username":     {"EN": "Username", "FR": "Pseudonyme"},
	"emailaddress": {"EN": "Email Address", "FR": "Adresse Email"},
	"message":      {"EN": "Message", "FR": "Message"},
	"search":       {"EN": "Search", "FR": "Recherche"},
}

func Translate(base string) string {
	val := translation[base]["FR"]
	if val != "" {
		return val
	}
	return "TRANSLATION ERROR"
}
