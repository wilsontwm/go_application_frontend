package utils

import (
	"github.com/gorilla/sessions"
	"net/http"
)

// Set the error/success flash message depends on the success state of the response
func SetErrorSuccessFlash(session *sessions.Session, w http.ResponseWriter, r *http.Request, resp map[string]interface{}) {
	// Set flash
	var messages []string

	if resp["errors"] != nil {
		messages = resp["errors"].([]string)
	} else {
		msg := resp["message"].(string)
		messages = append(messages, msg)
	}

	var tag string
	if resp["success"].(bool) {
		tag = "success"
	} else {
		tag = "errors"
	}

	for _, message := range messages {
		session.AddFlash(message, tag)
	}

	session.Save(r, w)
}
