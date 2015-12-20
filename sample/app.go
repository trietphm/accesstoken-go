// Sample app that is the equivalent to the official Python test app from Twilio

package main

import (
	"encoding/json"
	"fmt"
	"github.com/corytodd/accesstoken-go"
	"net/http"
	"os"
)

// Handles GET request for token
func token(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	// get credentials for environment variables
	accountSid := os.Getenv("TW_ACCOUNT_SID")
	apiKey := os.Getenv("TW_API_KEY")
	apiSecret := os.Getenv("TW_API_SECRET")

	// Create an Access Token
	myToken := accesstoken.New(accountSid, apiKey, apiSecret)

	// Set the Identity of this token
	id := "gotwilio.sample"
	myToken.Identity = id

	// Grant access to Conversations
	grant := accesstoken.NewConversationsGrant(os.Getenv("TW_VIDEO_SID"))
	myToken.AddGrant(grant)

	signedJWT, err := myToken.ToJWT(accesstoken.DEFAULT_ALGORITHM)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	resp := &Response{
		"identity": id,
		"token":    signedJWT,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	// Return token info as JSON
	w.WriteHeader(200)
	w.Write(b)
}

//Convenience interface for printing anonymous JSON objects
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/token", token)
	http.ListenAndServe(":8080", nil)
}