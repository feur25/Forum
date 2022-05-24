package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func randomString(len int) string {

	bytes := make([]byte, len)

	for i := 0; i < len; i++ {
		bytes[i] = byte(randInt(65, 90))
	}

	return string(bytes)
}
func randInt(min int, max int) int {

	return min + rand.Intn(max-min)
}

func send_update_email(recipient string) string {
	generate_code := randomString(5)

	from := "goeasycode@gmail.com"
	to := []string{recipient}
	msg := []byte("To: " + recipient + "\r\n" +
		"From: goeasycode@gmail.com\r\n" +
		"Subject: Code De vérification \r\n" +
		"\r\n" +
		"Voici votre code de vérification " +
		"code :\r\n" + generate_code)

	status := sendEmail(from, to, msg)

	if status {
		fmt.Printf("Email sent successfully.\n")
	} else {
		fmt.Printf("Email sent failed.\n")
	}

	return string(generate_code)
}
func send_delete_email(recipient string) string {
	from := "goeasycode@gmail.com"
	generate_code := randomString(5)

	to := []string{recipient}
	msg := []byte("To: " + recipient + "\r\n" +
		"From: goeasycode@gmail.com\r\n" +
		"Subject: Suppression du compte \r\n" +
		"\r\n" +
		"Code pour supprimer votre compte, il ne pourra être récupéré par la suite, une suppression est définitive !  " +
		"code :\r\n" + generate_code)

	status := sendEmail(from, to, msg)

	if status {
		fmt.Printf("Email sent successfully.\n")
	} else {
		fmt.Printf("Email sent failed.\n")
	}

	return string(generate_code)
}

func sendEmail(from string, to []string, msg []byte) bool {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file\n")
	}

	auth := smtp.PlainAuth("", from, os.Getenv("PASSWORD"), os.Getenv("SMTP_HOST"))

	smtpAddress := fmt.Sprintf("%s:%v", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))

	err = smtp.SendMail(smtpAddress, auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
func ping(senderName, recipientMail, message string) {

	from := "goeasycode@gmail.com"
	to := []string{recipientMail}
	msg := []byte("To: " + recipientMail + "\r\n" +
		"From: goeasycode@gmail.com\r\n" +
		"Subject: Quelqu'un vous a envoyé un message ! \r\n" +
		"\r\n" +
		senderName + " vous a ping" +
		"\nMessage : \r\n" + message)

	status := sendEmail(from, to, msg)

	if status {
		fmt.Printf("Email sent successfully.\n")
	} else {
		fmt.Printf("Email sent failed.\n")
	}
}
