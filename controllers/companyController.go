package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	util "app_frontend/utils"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"	
	"github.com/gorilla/csrf"
)

var CompanyIndexPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	// Set the URL path
	restURL.Path = "/api/dashboard/company/index"
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
		
		var companies []interface{}
		_, hasData := resp["companies"]

		if hasData {
			companies = resp["companies"].([]interface{})
		} 

		data := map[string]interface{}{
			"title": "My Company",
			"appName": appName,
			"appVersion": appVersion,
			"name": name,
			"picture": picture,
			"year": year,
			"companies": companies,
			"createURL": "/dashboard/company/store",
			csrf.TemplateTag: csrf.TemplateField(r),
		}

		data, err = util.InitializePage(w, r, store, data)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	
		err = templates.ExecuteTemplate(w, "company_index_html", data)
	
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}	
	}
}

var CompanyCreateSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/dashboard/company/store"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")
	
	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	slug := strings.TrimSpace(r.Form.Get("slug"))
	description := strings.TrimSpace(r.Form.Get("description"))
	email := strings.TrimSpace(r.Form.Get("email"))	
	phone := strings.TrimSpace(r.Form.Get("phone"))	
	fax := strings.TrimSpace(r.Form.Get("fax"))	
	address := strings.TrimSpace(r.Form.Get("address"))

	// Set the input data
	jsonData := map[string]interface{}{
		"name": name,
		"slug": slug,
		"description": description,
		"email": email,
		"phone": phone,
		"fax": fax,
		"address": address,
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

var CompanyShowPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	session, err := util.GetSession(store, w, r)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/show"
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
		
		if(resp["success"].(bool)) {
			var company map[string]interface{}
			_, hasData := resp["data"]
	
			if hasData {
				company = resp["data"].(map[string]interface{})
			} 
			
			data := map[string]interface{}{
				"title": company["Name"],
				"appName": appName,
				"appVersion": appVersion,
				"name": name,
				"picture": picture,
				"year": year,
				"company": company,
				"companyURL": appURL + "/dashboard/comp/" + company["Slug"].(string),
				csrf.TemplateTag: csrf.TemplateField(r),
			}
			
			data, err = util.InitializePage(w, r, store, data)
			
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		
			err = templates.ExecuteTemplate(w, "company_show_html", data)
		
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}	
		} else {
			util.SetErrorSuccessFlash(session, w, r, resp)
			// Redirect back to the previous page
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		}
	}
}

var CompanyShowJson = func(w http.ResponseWriter, r *http.Request){
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/show"
	urlStr := restURL.String()

	// Get the info for edit profile
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

var CompanyEditSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/update"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")
	
	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	slug := strings.TrimSpace(r.Form.Get("slug"))
	description := strings.TrimSpace(r.Form.Get("description"))
	email := strings.TrimSpace(r.Form.Get("email"))	
	phone := strings.TrimSpace(r.Form.Get("phone"))	
	fax := strings.TrimSpace(r.Form.Get("fax"))	
	address := strings.TrimSpace(r.Form.Get("address"))

	// Set the input data
	jsonData := map[string]interface{}{
		"name": name,
		"slug": slug,
		"description": description,
		"email": email,
		"phone": phone,
		"fax": fax,
		"address": address,
	}
	
	response, err := util.SendAuthenticatedRequest(urlStr, "PATCH", auth, jsonData)
	
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

var CompanyDeleteSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/delete"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

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
		
		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, "/dashboard/company", http.StatusFound)
	}
}

var CompanyGetUniqueSlugJson = func(w http.ResponseWriter, r *http.Request){
	var resp map[string]interface{}
	compQuery, ok := r.URL.Query()["comp"]
	companyId := ""
	if ok && len(compQuery[0]) >= 1 {
		companyId = compQuery[0]
	}

	slugQuery, ok := r.URL.Query()["slug"]
	slug := ""
	if ok && len(slugQuery[0]) >= 1 {
		slug = slugQuery[0]
	}

	// Set the URL path
	restURL.Path = "/api/dashboard/company/getUniqueSlug"
	queryString := restURL.Query()
	queryString.Set("comp", companyId)
	queryString.Set("slug", slug)
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