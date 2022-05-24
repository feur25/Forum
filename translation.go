package main

var translation = map[string]map[string]string{
	"accept":                   {"EN": "Accept", "FR": "Accepter"},
	"deny":                     {"EN": "Deny", "FR": "Refuser"},
	"send-message":             {"EN": "Message", "FR": "Envoyer un Message"},
	"delete":                   {"EN": "Delete", "FR": "Supprimer"},
	"send":                     {"EN": "Send", "FR": "Envoyer"},
	"welcome":                  {"EN": "Welcome", "FR": "Bienvenue"},
	"cancel":                   {"EN": "Cancel", "FR": "Annuler"},
	"login":                    {"EN": "Log in", "FR": "Se Connecter"},
	"logout":                   {"EN": "Log out", "FR": "Se Déconnecter"},
	"register":                 {"EN": "Register", "FR": "S'inscrire"},
	"username":                 {"EN": "Username", "FR": "Pseudonyme"},
	"update":                   {"EN": "update", "FR": "Mise a Jour"},
	"password":                 {"EN": "Password", "FR": "Mot de Passe"},
	"address":                  {"EN": "Address", "FR": "Adresse"},
	"email-address":            {"EN": "Email Address", "FR": "Adresse Email"},
	"phone-number":             {"EN": "Phone Number", "FR": "Numéro de Téléphone"},
	"first-name":               {"EN": "First Name", "FR": "Prénom"},
	"last-name":                {"EN": "Last Name", "FR": "Nom"},
	"message":                  {"EN": "Message", "FR": "Message"},
	"search":                   {"EN": "Search", "FR": "Recherche"},
	"send-friend-request":      {"EN": "Send a Friend Request", "FR": "Envoyez une Demande d'Ami"},
	"friends":                  {"EN": "Friends", "FR": "Amis"},
	"sent-friend-requests":     {"EN": "Sent Friend Requests", "FR": "Demandes d'Amis Envoyées"},
	"received-friend-requests": {"EN": "Search", "FR": "Demandes d'Amis Reçues"},
	"date":                     {"EN": "Creation date", "FR": "Date de création du compte"},
	"topic":                    {"EN": "Post New Topic", "FR": "Publier un Nouveau Topic"},
}

func Translate(base string) string {
	val := translation[base][data.Language]
	if val != "" {
		return val
	}
	return "TRANSLATION ERROR"
}
