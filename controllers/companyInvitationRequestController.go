package controllers

import (
	util "app_frontend/utils"
	"encoding/json"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
		"emails":  emails,
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
		responseBody, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(responseBody, &resp)

		if resp["success"].(bool) {
			// Send email to those that are invited
			company := resp["company"].(string)
			invitedEmails := resp["emails"].([]interface{})
			for _, invitedEmail := range invitedEmails {
				invitedEmailData := invitedEmail.(map[string]interface{})
				email := invitedEmailData["Email"].(string)
				link := appURL + "/dashboard/company/" + invitedEmailData["ID"].(string) + "/join"

				mailData := map[string]string{
					"appName":  appName,
					"joinLink": link,
					"company":  company,
					"message":  message,
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
		json.Unmarshal(responseBody, &resp)
	}

	util.Respond(w, resp)
}

// Resend the email to remind of the invitation requests
var CompanyInvitationResendSubmit = func(w http.ResponseWriter, r *http.Request) {
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
		json.Unmarshal(responseBody, &resp)

		// Send the email if it is a valid invitation
		_, hasData := resp["data"]
		_, hasCompanyData := resp["company"]
		if hasData && hasCompanyData && resp["success"].(bool) {
			invitation := resp["data"].(map[string]interface{})
			company := resp["company"].(map[string]interface{})
			link := appURL + "/dashboard/invite/incoming"
			email := invitation["Email"].(string)
			message := invitation["Message"].(string)
			companyName := company["Name"].(string)

			mailData := map[string]string{
				"appName":  appName,
				"joinLink": link,
				"company":  companyName,
				"message":  message,
			}
			subject := appName + " - You are invited!"

			r := util.NewRequest([]string{email}, subject)
			go r.Send("views/mail/invitation.html", mailData)
		}
	}

	util.Respond(w, resp)
}

// Resend the email to remind of the invitation requests
var CompanyInvitationResendMultipleSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Get the input data from the form
	r.ParseForm()
	invitationIdsString := strings.TrimSpace(r.Form.Get("invitationIds"))
	invitationIds := strings.Split(invitationIdsString, ",")

	auth := ReadEncodedCookieHandler(w, r, "auth")
	jsonData := make(map[string]interface{})

	// Loop through the invitation IDs and resend invitation emails
	// Create channel to receive the result
	const noOfInvitationWorkers int = 10 // Have 10 goroutines to get the emails
	invitationJobs := make(chan string, len(invitationIds))
	invitationEmails := make(chan string, len(invitationIds))

	for w := 1; w <= noOfInvitationWorkers; w++ {
		go func(invitationJobs <-chan string, results chan<- string) {
			for invitationInput := range invitationJobs {
				// Set the URL path
				restURL.Path = "/api/dashboard/company/" + companyId + "/invite/" + invitationInput
				urlStr := restURL.String()

				// Get the company invitation request
				response, err := util.SendAuthenticatedRequest(urlStr, "GET", auth, jsonData)
				email := ""
				if err == nil {
					responseBody, _ := ioutil.ReadAll(response.Body)

					// Parse it to json data
					json.Unmarshal(responseBody, &resp)

					// Send the email if it is a valid invitation
					_, hasData := resp["data"]
					_, hasCompanyData := resp["company"]
					if hasData && hasCompanyData && resp["success"].(bool) {
						invitation := resp["data"].(map[string]interface{})
						company := resp["company"].(map[string]interface{})
						link := appURL + "/dashboard/invite/incoming"
						message := invitation["Message"].(string)
						companyName := company["Name"].(string)

						// Only resend invitation if the status is still pending to avoid spamming
						if int(invitation["Status"].(float64)) == 0 {
							email = invitation["Email"].(string)
							mailData := map[string]string{
								"appName":  appName,
								"joinLink": link,
								"company":  companyName,
								"message":  message,
							}
							subject := appName + " - You are invited!"

							r := util.NewRequest([]string{email}, subject)
							go r.Send("views/mail/invitation.html", mailData)
						}
					}
				}

				results <- email
			}
		}(invitationJobs, invitationEmails)
	}

	// Loop through the emails to check if the email can be invited
	for _, invitationId := range invitationIds {
		// Send the email to the email jobs
		invitationJobs <- invitationId
	}
	close(invitationJobs)

	// Gather the result
	var successfulEmails []string
	for i := 0; i < len(invitationIds); i++ {
		successfulEmail := <-invitationEmails
		if successfulEmail != "" {
			successfulEmails = append(successfulEmails, successfulEmail)
		}
	}

	var errors []string
	if len(successfulEmails) > 0 {
		emails := strings.Join(successfulEmails, ", ")
		resp = util.Message(true, http.StatusOK, "You have successfully resend invitation to "+emails+".", errors)
	} else {
		resp = util.Message(false, http.StatusOK, "Something wrong has occured. Please ensure you have selected emails to resend invitation.", errors)
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
		responseBody, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(responseBody, &resp)

		util.Respond(w, resp)
	}
}

// Delete the invitation requests in bulk
var CompanyInvitationDeleteMultipleSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Get the input data from the form
	r.ParseForm()
	invitationIdsString := strings.TrimSpace(r.Form.Get("invitationIds"))
	invitationIds := strings.Split(invitationIdsString, ",")

	auth := ReadEncodedCookieHandler(w, r, "auth")
	jsonData := make(map[string]interface{})

	// Loop through the invitation IDs and delete invitation requests
	// Create channel to receive the result
	const noOfInvitationDeleteWorkers int = 10 // Have 10 goroutines to get the emails
	invitationJobs := make(chan string, len(invitationIds))
	invitationDeletedIds := make(chan string, len(invitationIds))

	for w := 1; w <= noOfInvitationDeleteWorkers; w++ {
		go func(invitationJobs <-chan string, results chan<- string) {
			for invitationInput := range invitationJobs {
				// Set the URL path
				restURL.Path = "/api/dashboard/company/" + companyId + "/invite/" + invitationInput + "/delete"
				urlStr := restURL.String()

				// Get the company invitation request
				response, err := util.SendAuthenticatedRequest(urlStr, "DELETE", auth, jsonData)
				deletedID := ""
				if err == nil {
					responseBody, _ := ioutil.ReadAll(response.Body)

					// Parse it to json data
					json.Unmarshal(responseBody, &resp)

					// Get the result from the response
					if _, ok := resp["success"]; ok && resp["success"].(bool) {
						deletedID = invitationInput
					}
				}

				results <- deletedID
			}
		}(invitationJobs, invitationDeletedIds)
	}

	// Loop through the emails to check if which invitation
	for _, invitationId := range invitationIds {
		// Send the email to the email jobs
		invitationJobs <- invitationId
	}
	close(invitationJobs)

	// Gather the result
	var deletedIDs []string
	for i := 0; i < len(invitationIds); i++ {
		deletedID := <-invitationDeletedIds
		if deletedID != "" {
			deletedIDs = append(deletedIDs, deletedID)
		}
	}

	var errors []string
	if len(deletedIDs) > 0 {
		resp = util.Message(true, http.StatusOK, "You have successfully deleted "+string(len(deletedIDs))+" invitation(s).", errors)
		resp["data"] = deletedIDs
	} else {
		resp = util.Message(false, http.StatusOK, "Something wrong has occured. Please ensure you have selected some invitations.", errors)
	}

	util.Respond(w, resp)
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
			"title":          "My Invites",
			"appName":        appName,
			"appVersion":     appVersion,
			"name":           name,
			"picture":        picture,
			"year":           year,
			"invitations":    invitations,
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
		responseBody, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(responseBody, &resp)

		if resp["success"].(bool) {
			data := resp["data"].(map[string]interface{})

			// Set the company into the redis when user accepts the invitation request
			if data["Status"].(float64) == 1.0 {
				// Get all the companies the user belongs to
				type Company struct {
					ID   string
					Name string
				}

				id := ReadCookieHandler(w, r, "id")
				comp := resp["company"].(map[string]interface{})
				company := Company{}
				compJsonBody, _ := json.Marshal(comp)
				json.Unmarshal(compJsonBody, &company)
				redisdata, _ := json.Marshal(&company)
				util.RedisRPush("user:"+id+";companies:", []byte(string(redisdata)))
			}
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}
