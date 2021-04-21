package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// saves the token
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// for labels
	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}

	// For message
	// p, err := srv.Users.Messages.List(user).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Messages: %v", err)
	// }
	// if len(p.Messages) == 0 {
	// 	fmt.Println("No Messages found.")
	// 	return
	// }
	// fmt.Println("Messages:")
	// for i, n := range p.Messages {
	// 	msg, err := srv.Users.Messages.Get("me", n.Id).Do()
	// 	if err != nil {
	// 		log.Fatalf("Unable to retrieve message %v: %v", n.Id, err)
	// 	}
	// 	fmt.Println(i, msg.Id, msg.Snippet)

	// }

	// for Drafts
	// t, err := srv.Users.Drafts.List(user).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Drafts: %v", err)
	// }
	// if len(t.Drafts) == 0 {
	// 	fmt.Println("No Messages found.")
	// 	return
	// }
	// fmt.Println("Drafts:")
	// for i, m := range t.Drafts {
	// 	draft, err := srv.Users.Drafts.Get("me", m.Id).Do()
	// 	if err != nil {
	// 		log.Fatalf("Unable to retrieve draft %v: %v", m.Id, err)
	// 	}
	// 	fmt.Println(i, draft)
	// 	// fmt.Println(i, m.Message.Id, m.ForceSendFields)
	// }

}
