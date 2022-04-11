package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

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

func send_reset_mail() {
	from := "go.easy.bot@gmail.com"
	rand.Seed(time.Now().UTC().UnixNano())
	generate_code := randomString(5)
	to := []string{data.Auth.email}
	msg := []byte("To: " + data.Auth.email + "\r\n" +
		"From: go.easy.bot@gmail.com\r\n" +
		"Subject: Code De vérification \r\n" +
		"\r\n" +
		"Voici votre code de vérification " +
		"code :\r\n" + generate_code)
	data.code = string(generate_code)

	status := sendEmail(from, to, msg)

	if status {
		fmt.Printf("Email sent successfully.\n")
	} else {
		fmt.Printf("Email sent failed.\n")
	}
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
