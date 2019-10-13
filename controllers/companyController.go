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

// Show a list of company that the user belongs to
var CompanyIndexPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	// Set the URL path
	restURL.Path = "/api/dashboard/company"
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
		json.Unmarshal(responseBody, &resp)

		var companies []interface{}
		_, hasData := resp["companies"]

		if hasData {
			companies = resp["companies"].([]interface{})
		}

		data := map[string]interface{}{
			"title":          "My Company",
			"appName":        appName,
			"appVersion":     appVersion,
			"name":           name,
			"picture":        picture,
			"year":           year,
			"companies":      companies,
			"createURL":      "/dashboard/company/store",
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

// Create a new company
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
		"name":        name,
		"slug":        slug,
		"description": description,
		"email":       email,
		"phone":       phone,
		"fax":         fax,
		"address":     address,
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
		json.Unmarshal(data, &resp)

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// Show the details of the company specified
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
		json.Unmarshal(responseBody, &resp)

		if resp["success"].(bool) {
			company := make(map[string]interface{})
			isAdmin := false

			if _, ok := resp["data"]; ok {
				company = resp["data"].(map[string]interface{})
			}

			if _, ok := resp["isAdmin"]; ok {
				isAdmin = resp["isAdmin"].(bool)
			}

			data := map[string]interface{}{
				"title":          company["Name"],
				"appName":        appName,
				"appVersion":     appVersion,
				"name":           name,
				"picture":        picture,
				"year":           year,
				"company":        company,
				"companyURL":     appURL + "/dashboard/comp/" + company["Slug"].(string),
				"isAdmin":        isAdmin,
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

// Get the company detail in json
var CompanyShowJson = func(w http.ResponseWriter, r *http.Request) {
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
		json.Unmarshal(responseBody, &resp)
	}

	util.Respond(w, resp)
}

// Edit the company
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
		"name":        name,
		"slug":        slug,
		"description": description,
		"email":       email,
		"phone":       phone,
		"fax":         fax,
		"address":     address,
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
		json.Unmarshal(data, &resp)

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// Delete the company
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
		json.Unmarshal(data, &resp)

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, "/dashboard/company", http.StatusFound)
	}
}

// Get a unique slug/URL for the company via AJAX
var CompanyGetUniqueSlugJson = func(w http.ResponseWriter, r *http.Request) {
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
		json.Unmarshal(responseBody, &resp)
	}

	util.Respond(w, resp)
}

// Get a list of company users
var CompanyUsersListJson = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/users"
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

		if _, ok := resp["data"]; ok {
			if datas, ok := resp["data"].([]interface{}); ok {
				for _, data := range datas {
					if data, ok := data.(map[string]interface{}); ok {
						// Check if the profile picture is set, else set a default picture
						if data["profilePicture"] == "" {
							data["profilePicture"] = defaultProfilePic
						}
					}
				}

			}
		}
	}

	util.Respond(w, resp)
}

// Get a list of company users by search
var CompanyUsersSearchJson = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	var errors []string
	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

	if companyId == "" {
		resp := util.Message(false, http.StatusOK, "Please select a company first.", errors)
		util.Respond(w, resp)
		return
	}

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/users/search"
	queryString := restURL.Query()
	searchQuery, ok := r.URL.Query()["query"]
	if ok && len(searchQuery[0]) >= 1 {
		queryString.Set("query", searchQuery[0])
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

		// Set the redis
		if resp["success"].(bool) {
			if _, ok := resp["data"]; ok {
				if datas, ok := resp["data"].([]interface{}); ok {
					for _, data := range datas {
						if data, ok := data.(map[string]interface{}); ok {
							// Check if the profile picture is set, else set a default picture
							if data["profilePicture"] == "" {
								data["profilePicture"] = defaultProfilePic
							}
						}
					}

				}
			}
		}
	}

	util.Respond(w, resp)
}

// Get a list of company users by search (page)
var CompanyUsersSearchPage = func(w http.ResponseWriter, r *http.Request) {
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	data := map[string]interface{}{
		"title":          "Search",
		"appName":        appName,
		"appVersion":     appVersion,
		"name":           name,
		"picture":        picture,
		"year":           year,
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	data, err := util.InitializePage(w, r, store, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "profile_search_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Select a company as current active one
var CompanyVisitSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/visit"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Set the input data
	jsonData := make(map[string]interface{})

	response, err := util.SendAuthenticatedRequest(urlStr, "PATCH", auth, jsonData)

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

		// Set the redis
		if resp["success"].(bool) {
			type Company struct {
				ID   string
				Name string
			}

			id := ReadCookieHandler(w, r, "id")
			// Set selected company into redis
			selectedCompany := Company{}
			compJsonBody, _ := json.Marshal(resp["selectedCompany"].(map[string]interface{}))
			json.Unmarshal(compJsonBody, &selectedCompany)
			redisdata, _ := json.Marshal(&selectedCompany)
			util.RedisSet("user:"+id+";selectedcompany:", []byte(string(redisdata)))
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the dashboard page
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}
