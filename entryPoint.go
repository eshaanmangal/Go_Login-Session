package main

import (
	"encoding/json"
	"fmt"
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

// Welcome func
func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	response, valid := Database.GetValue(sessionToken)
	if !valid {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// if response == nil {
	// 	// If the session token is not present in cache, return an unauthorized error
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
}

func main() {
	// Database.SaveKeyValue("temp@email.com", "hahaha", 0)
	http.HandleFunc("/signin", signin)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
