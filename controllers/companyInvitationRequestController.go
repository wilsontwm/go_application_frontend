package controllers

import (
	"net/http"
	"time"
	"github.com/gorilla/mux"
	util "app_frontend/utils"
	"strings"
	"io/ioutil"
	"encoding/json"	
	"github.com/gorilla/csrf"
)

// Send invitation emails to join company
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
	message := strings.TrimSpace(r.Form.Get("message"))
	var emails []string
	json.Unmarshal([]byte(emailsString), &emails)

	jsonData := map[string]interface{}{
		"emails": emails,
		"message": message,
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
					"message": message,
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

// Get a list of company invitation requests
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

// Resend the email to remind of the invitation requests
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
			message := invitation["Message"].(string)
			companyName := company["Name"].(string)

			mailData := map[string]string{
				"appName": appName,
				"joinLink": link,
				"company": companyName,
				"message": message,
			}
			subject := appName + " - You are invited!"

			r := util.NewRequest([]string{email}, subject)
			go r.Send("views/mail/invitation.html", mailData)
		}
	}
	
	util.Respond(w, resp)
}

// Delete the company invitation request
var CompanyInvitationDeleteSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]
	invitationId := vars["invitationID"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/invite/" + invitationId + "/delete"
	urlStr := restURL.String()

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Set the input data
	jsonData := make(map[string]interface{})
	
	response, err := util.SendAuthenticatedRequest(urlStr, "DELETE", auth, jsonData)

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
		
		util.Respond(w, resp)
	}
}

// User gets all the invitation requests from all the companies
var IndexInvitationFromCompany = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	// Set the URL path
	restURL.Path = "/api/dashboard/invite/incoming"
	urlStr := restURL.String()

	// Get the info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")
	jsonData := make(map[string]interface{})
	response, err := util.SendAuthenticatedRequest(urlStr, "GET", auth, jsonData)
	
	// Check if response is unauthorized
	if !CheckAuthenticatedRequest(w, r, response.StatusCode) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseBody, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(responseBody)), &resp)
		
		var invitations []interface{}
		_, hasData := resp["data"]

		if hasData {
			invitations = resp["data"].([]interface{})
		} 

		data := map[string]interface{}{
			"title": "My Invites",
			"appName": appName,
			"appVersion": appVersion,
			"name": name,
			"picture": picture,
			"year": year,
			"invitations": invitations,
			csrf.TemplateTag: csrf.TemplateField(r),
		}

		data, err = util.InitializePage(w, r, store, data)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	
		err = templates.ExecuteTemplate(w, "company_invitations_index_html", data)
	
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}	
	}
}

// Respond to the company invitation request
var RespondCompanyInvitationRequestSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	invitationId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/invite/incoming/" + invitationId + "/respond"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	responseString := strings.TrimSpace(r.Form.Get("response"))
	isJoin := true
	if responseString == "decline" {
		isJoin = false
	}
	// Set the input data
	jsonData := map[string]interface{}{
		"is_join": isJoin,
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

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}