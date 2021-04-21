# gmail OAuth2 in Go

This is implement your own authentication.

Firstly create Google OAuth2 keys
Go to https://console.cloud.google.com/apis/dashboard
1. Create new project
2. Enable API & Services 
3. Welcome to the API Library
4. Choose Gmail API and Enable
5. Click "Create credentials"
6. Then Choose Consent Screen
7. Add authorized redirect URL, in our case it will be localhost:8080/callback
8. Get client id and client secret
9. Save it in a safe place

We'll do everything in main.go file, and register 3 URL handlers:

1. /
2. /login
3. /callback

Initial handlers and OAuth2 config
go get golang.org/x/oauth2
go get cloud.google.com/go/compute/metadata

Test it
go run main.go


User information, messages, labels, drafts can be easily viewed using GoLang through Gmail API.
