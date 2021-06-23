package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antonholmquist/jason"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// NOTE : we don't want to visit CSRF URL to get the authorization code
// and paste into the terminal each time we want to send an email
// therefore we will retrieve a token for our client, save the token into a file
// you will be prompted to visit a link in your browser for authorization code only ONCE
// and subsequent execution of the program will not prompt you for authorization code again
// until the token expires.

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// dir, err := os.Getwd()
	// fmt.Println(dir)
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	// tokenCacheDir := filepath.Join(dir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-sendemail.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func chunkSplit(body string, limit int, end string) string {

	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result string = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit

		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		// charSlice
		if len(charSlice) < limit {
			limit = len(charSlice)
		}

	}

	return result

}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func main() {
	ctx := context.Background()

	// process the credential file
	credential, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// In order for POST upload attachment to work
	// You need to authorize the Gmail API v1 scope
	// at https://developers.google.com/oauthplayground/
	// otherwise you will get Authorization error in the API JSON reply

	// Use MailGoogleComScope for this example. Because of Upload Attachments
	// and Draft creation.  see https://developers.google.com/gmail/api/auth/scopes

	// See the rest at https://godoc.org/google.golang.org/api/gmail/v1#pkg-constants

	config, err := google.ConfigFromJSON(credential, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)

	// initiate a new gmail client service
	gmailClientService, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to initiate new gmail client: %v", err)
	}

	// get our token
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	token, err := tokenFromFile(cacheFile)
	if err != nil {
		log.Fatalf("Unable to get token from file. %v", err)
	}

	// read file for attachment purpose
	fileName := "img.pdf"
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Unable to read file for attachment: %v", err)
	}

	fileMIMEType := http.DetectContentType(fileBytes)

	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	fileSize := len(string(fileBytes))

	content := `Hello!
                     this is a draft message with attachment.
                     Bye`

	userID := "me"

	postURL := "https://www.googleapis.com/upload/gmail/v1/users/" + userID + "/drafts?uploadType=media"

	// extract auth or access token from Token file
	// see https://godoc.org/golang.org/x/oauth2#Token
	authToken := token.AccessToken

	boundary := randStr(32, "alphanum")

	emailMessageData := []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"from: " + "i.joni40@gmail.com" + "\n" +
		"to: " + "asgor.ice@gmail.com" + "\n" +
		"subject: " + "upload attachment to draft and send" + "\n\n" +

		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		content + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName + string('"') + " \n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: base64\n" +
		"Content-Disposition: attachment; filename=" + string('"') + fileName + string('"') + " \n\n" +
		chunkSplit(fileData, 76, "\n") +
		"--" + boundary + "--")

	// convert []byte to io.Reader type with bytes.NewBuffer
	// request, _ := http.NewRequest("POST", postURL, bytes.NewBuffer(emailMessageData))

	// or

	// convert []byte to io.Reader type with strings.NewReader
	// see https://www.socketloop.com/tutorials/golang-convert-cast-byte-to-io-reader-type
	request, _ := http.NewRequest("POST", postURL, strings.NewReader(string(emailMessageData)))

	// see https://www.socketloop.com/tutorials/golang-post-data-with-url-values
	request.Header.Add("Host", "www.googleapis.com")
	request.Header.Add("Content-Type", "message/rfc822")
	request.Header.Add("Content-Length", strconv.Itoa(fileSize))
	request.Header.Add("Authorization", "Bearer "+authToken)

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Unable to be post to Google API: %v", err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Unable to read Google API response: %v", err)
	}

	// output the response from GMail API
	fmt.Println(string(body))

	// we need to extract the draft message ID to execute Send command
	jsonAPIreply, _ := jason.NewObjectFromBytes(body)

	draftID, _ := jsonAPIreply.GetString("id")
	fmt.Println("Draft ID : ", draftID)

	// ----- comment the lines below and you will send that the draft will appear
	// ----- in your Gmail's Draft box instead of being send out

	//https://godoc.org/google.golang.org/api/gmail/v1#Draft
	var draft gmail.Draft
	draft.Id = draftID

	// send out our draft
	_, err = gmailClientService.Users.Drafts.Send(userID, &draft).Do()
	if err != nil {
		log.Fatalf("Unable to send message: %v", err)
	} else {
		log.Println("Draft email with ID " + draftID + " sent!")
	}

}
