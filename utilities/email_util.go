package utilities

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/raLaaaa/gorala/models"
	"github.com/raLaaaa/gorala/services"
)

type EmailUtil struct{}

const url string = "https://api.sendinblue.com/v3/smtp/email"

func (e EmailUtil) SendRegistryConfirmation(token models.ConfirmationToken, user models.User) {

	dbService := services.DatabaseService{}

	mailsToday, err := dbService.FindAllEmailsOfToday()

	if err != nil {
		panic(err)
	}

	if len(mailsToday) > 280 {
		fmt.Println("Too many mails have been sent")
		return
	}

	var jsonStr = []byte(`{  
		"sender":{  
		   "name":"no-reply@gorala.icu",
		   "email":"no-reply@gorala.icu"
		},
		"to":[  
		   {  
			  "email":"` + user.Email + `",
			  "name":"` + user.Email + `"
		   }
		],
		"subject":"Registration Confirmation",
		"htmlContent":"<html><head></head><body><p>Hello,</p>Please click the following link to complete your registration.</p><p>http://localhost:1323/register/confirm/` + token.Token + `</p><p>Best regards</p></body></html>"
	 }`)

	jsonObj := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", url, jsonObj)

	if err != nil {
		panic(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", os.Getenv("SIBKEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	mail := models.SentMail{
		Recipient: user.Email,
		Text:      "registry",
	}

	dbService.CreateSentMail(&mail)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func (e EmailUtil) SendResetPassword(token models.ResetToken, user models.User) {

	dbService := services.DatabaseService{}

	mailsToday, err := dbService.FindAllEmailsOfToday()

	if err != nil {
		panic(err)
	}

	if len(mailsToday) > 280 {
		fmt.Println("Too many mails have been sent")
		return
	}

	fmt.Println(len(mailsToday))

	var jsonStr = []byte(`{  
		"sender":{  
		   "name":"no-reply@gorala.icu",
		   "email":"no-reply@gorala.icu"
		},
		"to":[  
		   {  
			"email":"` + user.Email + `",
			"name":"` + user.Email + `"
		   }
		],
		"subject":"Reset Password",
		"htmlContent":"<html><head></head><body><p>Hello,</p>Please click the following link to reset your password.</p><p>http://localhost:1323/reset/` + token.Token + `</p><p>Best regards</p></body></html>"
	 }`)

	jsonObj := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", url, jsonObj)

	if err != nil {
		panic(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", os.Getenv("SIBKEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	mail := models.SentMail{
		Recipient: user.Email,
		Text:      "resetpw",
	}

	dbService.CreateSentMail(&mail)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
