# gmail OAuth2 in Go

Implement your own authentication.

Firstly create Google OAuth2 keys
# Go to https://console.cloud.google.com/apis/dashboard
1. Create new project
2. Enable API & Services 
3. Welcome to the API Library
4. Choose Gmail API and Enable
5. Click "Create credentials"
6. Then Choose Consent Screen
7. Add authorized redirect URL, in our case it will be localhost:8080/callback
8. Get client id and client secret
9. Save it in a safe place

# We've done everything in the main.go file:

1. /
2. /login
3. /callback

# Initial handlers and OAuth2 config
1. go get golang.org/x/oauth2/google
2. go get google.golang.org/api/gmail/v1


# Test and Run
 go run main.go


User information, messages, labels, drafts are easily accessible using Golang through Gmail API.
