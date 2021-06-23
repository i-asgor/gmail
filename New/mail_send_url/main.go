package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var google0authConfig = &oauth2.Config{
	ClientID:     "923806387165-5n065ckbcdvg3ht7plrp9m1q0sg0beda.apps.googleusercontent.com",
	ClientSecret: "n-bdTBwCHXllGv-Ehg8COVe0",
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/gmail.readonly",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/gmail.compose",
		"https://www.googleapis.com/auth/gmail.send"},
	RedirectURL: "http://localhost:8080/callback",
	Endpoint:    google.Endpoint,
}
var randomState = "random"

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/login">Google Log In</a></body></html>`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := google0authConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// conf, err := google.ConfigFromJSON(google0authConfig, gmail.GmailComposeScope)
	// if err != nil {
	// 	log.Printf("Error: %v", err)
	// }

	token, err := google0authConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	// fmt.Println(token)
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// resp, err := http.Get("https://gmail.googleapis.com/gmail/v1/users/i.joni40@gmail.com/messages/send?alt=json&access_token=" + token.AccessToken)
	// if err != nil {
	// 	fmt.Printf("could not create get request: %s\n", err.Error())
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// defer resp.Body.Close()
	// // fmt.Println(resp.Body)
	// content, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("could not parse response: %s\n", err.Error())
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// fmt.Fprintf(w, "Response: %s", content)

	// Create the *http.Client using the access token
	client := google0authConfig.Client(oauth2.NoContext, token)

	// Create a new gmail service using the client
	gmailService, err := gmail.New(client)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	// New message for our gmail service to send
	var message gmail.Message

	// Compose the message
	messageStr := []byte(
		"From: i.joni40@gmail.com\r\n" +
			"To: i.joni0640@gmail.com\r\n" +
			"Subject: My first Gmail API message\r\n\r\n" +
			"Message body goes here!")

	// Place messageStr into message.Raw in base64 encoded format
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	// Send the message
	_, err = gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}

	fmt.Fprintf(w, "Response: %s", messageStr)

}
