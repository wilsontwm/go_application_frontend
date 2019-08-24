package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	util "app_frontend/utils"
	"strings"
	"io/ioutil"
	"encoding/json"	
)

var CompanyInviteSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/invite"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	emailsString := strings.TrimSpace(r.Form.Get("emails"))
	var emails []string
	json.Unmarshal([]byte(emailsString), &emails)

	jsonData := map[string]interface{}{
		"emails": emails,
	}
	
	response, err := util.SendAuthenticatedRequest(urlStr, "POST", auth, jsonData)

	// Check if response is unauthorized
	if !CheckAuthenticatedRequest(w, r, response.StatusCode) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)	

		if(resp["success"].(bool)) {
			// Send email to those that are invited
			company := resp["company"].(string)
			invitedEmails := resp["emails"].([]interface{})
			for _, invitedEmail := range invitedEmails {
				invitedEmailData := invitedEmail.(map[string]interface{})
				email := invitedEmailData["Email"].(string)
				link := appURL + "/dashboard/company/" + invitedEmailData["ID"].(string) + "/join"

				mailData := map[string]string{
					"appName": appName,
					"joinLink": link,
					"company": company,
				}
				subject := appName + " - You are invited!"

				r := util.NewRequest([]string{email}, subject)
				go r.Send("views/mail/invitation.html", mailData)
			}
		} 

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

var CompanyInviteListJson = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/invite/list"
	queryString := restURL.Query()
	pageQuery, ok := r.URL.Query()["page"]
	if ok && len(pageQuery[0]) >= 1 {		
		queryString.Set("page", pageQuery[0])
	}
	
	restURL.RawQuery = queryString.Encode()
	urlStr := restURL.String()

	// Check if the URL is unique
	auth := ReadEncodedCookieHandler(w, r, "auth")
	jsonData := make(map[string]interface{})
	response, err := util.SendAuthenticatedRequest(urlStr, "GET", auth, jsonData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseBody, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(responseBody)), &resp)
	}
	
	util.Respond(w, resp)
}

var CompanyInvitationResendSubmit = func(w http.ResponseWriter, r *http.Request){
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]
	invitationId := vars["invitationID"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/invite/" + invitationId
	urlStr := restURL.String()

	// Get the company invitation request
	auth := ReadEncodedCookieHandler(w, r, "auth")
	jsonData := make(map[string]interface{})
	response, err := util.SendAuthenticatedRequest(urlStr, "GET", auth, jsonData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseBody, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(responseBody)), &resp)
		
		// Send the email if it is a valid invitation
		_, hasData := resp["data"]
		_, hasCompanyData := resp["company"]
		if(hasData && hasCompanyData && resp["success"].(bool)) {
			invitation := resp["data"].(map[string]interface{})
			company := resp["company"].(map[string]interface{})
			link := appURL + "/dashboard/company/" + invitation["ID"].(string) + "/join"
			email := invitation["Email"].(string)
			companyName := company["Name"].(string)

			mailData := map[string]string{
				"appName": appName,
				"joinLink": link,
				"company": companyName,
			}
			subject := appName + " - You are invited!"

			r := util.NewRequest([]string{email}, subject)
			go r.Send("views/mail/invitation.html", mailData)
		}
	}
	
	util.Respond(w, resp)
}