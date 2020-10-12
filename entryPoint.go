package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/login-go/Database"
	uuid "github.com/satori/go.uuid"
)

// Credentials of users
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func signin(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := Database.GetValue(creds.Username)

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewV4().String()
	Database.SaveKeyValue(sessionToken, creds.Username, 120)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func main() {
	// Database.SaveKeyValue("temp@email.com", "hahaha", 0)
	http.HandleFunc("/signin", signin)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
