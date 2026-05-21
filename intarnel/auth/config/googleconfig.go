package gconfig

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


var GoogleConfig *oauth2.Config

func Intauth() {

	GoogleConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"),
		RedirectURL:  "http://localhost:8080/callbackfromgoogle",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

}
