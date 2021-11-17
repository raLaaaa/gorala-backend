package utilities

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/raLaaaa/gorala/models"
)

type EmailUtil struct{}

const url string = "https://api.sendinblue.com/v3/smtp/email"

func (e EmailUtil) SendRegistryConfirmation(token models.ConfirmationToken) {

	var jsonStr = []byte(`{  
		"sender":{  
		   "name":"Sender Alex",
		   "email":"senderalex@example.com"
		},
		"to":[  
		   {  
			  "email":"sageinenderlitterist@gmail.com",
			  "name":"John Doe"
		   }
		],
		"subject":"Hello world",
		"htmlContent":"<html><head></head><body><p>Hello,</p>This is my first transactional email sent from Sendinblue.</p></body></html>"
	 }'`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func (e EmailUtil) SendResetPassword(token models.ResetToken) {

	var jsonStr = []byte(`{  
		"sender":{  
		   "name":"Sender Alex",
		   "email":"senderalex@example.com"
		},
		"to":[  
		   {  
			  "email":"sageinenderlitterist@gmail.com",
			  "name":"John Doe"
		   }
		],
		"subject":"Hello world",
		"htmlContent":"<html><head></head><body><p>Hello,</p>This is my first transactional email sent from Sendinblue.</p></body></html>"
	 }'`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
