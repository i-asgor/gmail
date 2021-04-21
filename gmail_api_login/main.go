package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var google0authConfig = &oauth2.Config{
	ClientID:     "923806387165-5n065ckbcdvg3ht7plrp9m1q0sg0beda.apps.googleusercontent.com",
	ClientSecret: "n-bdTBwCHXllGv-Ehg8COVe0",
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/gmail.readonly",
		"https://www.googleapis.com/auth/userinfo.email"},
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

	token, err := google0authConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	// fmt.Println(token)
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	userID := "me"
	// resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=" + token.AccessToken)
	resp, err := http.Get("https://gmail.googleapis.com/gmail/v1/users/" + userID + "/messages?alt=json&access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create get request: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()
	// fmt.Println(resp.Body)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not parse response: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Response: %s", content)

	// var gmData interface{}
	// json.Unmarshal(content, &gmData) //extracting the json file

	// fmt.Println(gmData)

	// taking required values

	// payload := gmData.(map[string]interface{})["payload"]
	// parts := payload.(map[string]interface{})["parts"]
	// labelIds := gmData.(map[string]interface{})["labelIds"]
	// headers := payload.(map[string]interface{})["headers"]
	// HistoryId := gmData.(map[string]interface{})["historyId"]
	// snippet := gmData.(map[string]interface{})["snippet"]
	// fmt.Fprintf(w, "HistoryId: %s  %s %s", HistoryId, "\n", snippet)
}
