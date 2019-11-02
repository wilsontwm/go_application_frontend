package controllers

import (
	util "app_frontend/utils"
	"encoding/json"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const PostStatusDraft = "Draft"
const PostStatusScheduled = "Scheduled"
const PostStatusPublished = "Published"

var PostStatusArray = [...]string{
	PostStatusDraft,
	PostStatusScheduled,
	PostStatusPublished,
}

type PostStatus struct {
	ID     int
	Status string
}

// Show create post page
var PostCreatePage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	session, err := util.GetSession(store, w, r)

	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

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

		var status []PostStatus
		for i, stat := range PostStatusArray {
			s := PostStatus{
				ID:     i,
				Status: stat,
			}
			status = append(status, s)
		}

		// Parse it to json data
		json.Unmarshal(responseBody, &resp)

		if resp["success"].(bool) {
			data := map[string]interface{}{
				"title":          "New post",
				"appName":        appName,
				"appVersion":     appVersion,
				"name":           name,
				"picture":        picture,
				"year":           year,
				"postStatus":     status,
				"url":            "/dashboard/post/store",
				"isEdit":         false,
				csrf.TemplateTag: csrf.TemplateField(r),
			}

			data, err = util.InitializePage(w, r, store, data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = templates.ExecuteTemplate(w, "post_create_html", data)

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

// Create a new post
var PostCreateSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/post/store"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	title := strings.TrimSpace(r.Form.Get("title"))
	content := strings.TrimSpace(r.Form.Get("content"))
	status, _ := strconv.ParseInt(r.Form.Get("status"), 10, 32)
	dateFormat := "02 Jan 2006 15:04 PM" // DD MMM YYYY h:mm A
	scheduledAt, _ := time.ParseInLocation(dateFormat, r.Form.Get("scheduled_at"), time.Now().Location())

	// Set the input data
	jsonData := map[string]interface{}{
		"title":        title,
		"content":      content,
		"status":       status,
		"scheduled_at": scheduledAt,
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

		// Redirect back to the profile page with the post listing if successful
		// Else redirect to previous page
		if resp["success"].(bool) {
			url := "/dashboard/user/" + userId
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		}
	}
}

// Show the details of the post specified
var PostListJson = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/post"
	restURL.RawQuery = ""
	queryString := restURL.Query()

	authorQuery, ok := r.URL.Query()["author"]
	if ok && len(authorQuery[0]) >= 1 {
		queryString.Set("author", authorQuery[0])
	}

	statusQuery, ok := r.URL.Query()["status"]
	if ok && len(statusQuery[0]) >= 1 {
		queryString.Set("status", statusQuery[0])
	}

	idQuery, ok := r.URL.Query()["lastID"]
	if ok && len(idQuery[0]) >= 1 {
		queryString.Set("lastID", idQuery[0])
	}

	lastUpdatedQuery, ok := r.URL.Query()["lastUpdated"]
	if ok && len(lastUpdatedQuery[0]) >= 1 {
		queryString.Set("lastUpdated", lastUpdatedQuery[0])
	}

	lastPublishedQuery, ok := r.URL.Query()["lastPublished"]
	if ok && len(lastPublishedQuery[0]) >= 1 {
		queryString.Set("lastPublished", lastPublishedQuery[0])
	}

	limitQuery, ok := r.URL.Query()["limit"]
	if ok && len(limitQuery[0]) >= 1 {
		queryString.Set("limit", limitQuery[0])
	}

	restURL.RawQuery = queryString.Encode()
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

	resp["defaultProfilePic"] = defaultProfilePic // default profile picture
	util.Respond(w, resp)
}

// Show the details of the post specified
var PostShowPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	session, err := util.GetSession(store, w, r)

	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

	// Get the ID of the post passed in via URL
	vars := mux.Vars(r)
	postId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/post/" + postId + "/show"
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
			var errors []string
			post := make(map[string]interface{})

			if _, ok := resp["data"]; ok {
				post = resp["data"].(map[string]interface{})
				author := post["Author"].(map[string]interface{})
				if author["profilePicture"] == nil || author["profilePicture"] == "" {
					author["profilePicture"] = defaultProfilePic // default profile picture
				}

				// If the post is not published, then do not show post unless it's the author
				if post["PublishedAt"] == nil && post["AuthorID"] != userId {

					resp = util.Message(false, http.StatusOK, "You are unauthorized to perform the action.", errors)
					util.SetErrorSuccessFlash(session, w, r, resp)
					// Redirect back to the previous page
					http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
					return
				}

				data := map[string]interface{}{
					"title":          post["Title"],
					"appName":        appName,
					"appVersion":     appVersion,
					"name":           name,
					"picture":        picture,
					"year":           year,
					"post":           post,
					csrf.TemplateTag: csrf.TemplateField(r),
				}

				data, err = util.InitializePage(w, r, store, data)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				err = templates.ExecuteTemplate(w, "post_show_html", data)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				resp = util.Message(false, http.StatusOK, "You are unauthorized to perform the action.", errors)
				util.SetErrorSuccessFlash(session, w, r, resp)
				// Redirect back to the previous page
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			}
		} else {
			util.SetErrorSuccessFlash(session, w, r, resp)
			// Redirect back to the previous page
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		}
	}
}

// Show edit post page
var PostEditPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	session, err := util.GetSession(store, w, r)

	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

	// Get the ID of the post passed in via URL
	vars := mux.Vars(r)
	postId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/post/" + postId + "/show"
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
			post := make(map[string]interface{})

			if _, ok := resp["data"]; ok {
				post = resp["data"].(map[string]interface{})

				// If the post does not belong to the author
				if post["AuthorID"] != userId {
					var errors []string
					resp = util.Message(false, http.StatusOK, "You are unauthorized to perform the action.", errors)
					util.SetErrorSuccessFlash(session, w, r, resp)
					// Redirect back to the previous page
					http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
					return
				}
			}

			data := map[string]interface{}{
				"title":          post["Title"],
				"appName":        appName,
				"appVersion":     appVersion,
				"name":           name,
				"picture":        picture,
				"year":           year,
				"post":           post,
				"postStatus":     resp["postStatus"],
				"url":            "/dashboard/post/" + post["ID"].(string) + "/update",
				"deleteUrl":      "/dashboard/post/" + post["ID"].(string) + "/delete",
				"isEdit":         true,
				csrf.TemplateTag: csrf.TemplateField(r),
			}

			data, err = util.InitializePage(w, r, store, data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = templates.ExecuteTemplate(w, "post_edit_html", data)

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

// Edit existing post
var PostEditSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

	// Set the URL path
	// Get the ID of the post passed in via URL
	vars := mux.Vars(r)
	postId, _ := vars["id"]
	restURL.Path = "/api/dashboard/company/" + companyId + "/post/" + postId + "/update"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	title := strings.TrimSpace(r.Form.Get("title"))
	content := strings.TrimSpace(r.Form.Get("content"))
	status, _ := strconv.ParseInt(r.Form.Get("status"), 10, 32)
	dateFormat := "02 Jan 2006 15:04 PM" // DD MMM YYYY h:mm A
	scheduledAt, _ := time.ParseInLocation(dateFormat, r.Form.Get("scheduled_at"), time.Now().Location())

	// Set the input data
	jsonData := map[string]interface{}{
		"title":        title,
		"content":      content,
		"status":       status,
		"scheduled_at": scheduledAt,
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

		// Redirect back to the profile page with the post listing if successful
		// Else redirect to previous page
		if resp["success"].(bool) {
			url := "/dashboard/user/" + userId
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		}
	}
}

// Delete the post
var PostDeleteSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	userId := util.ReadCookieHandler(w, r, "id")
	companyId := util.GetActiveCompanyID(w, r, userId)

	// Get the ID of the post passed in via URL
	vars := mux.Vars(r)
	postId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/company/" + companyId + "/post/" + postId + "/delete"
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

		// Redirect back to the profile page with the post listing if successful
		// Else redirect to previous page
		if resp["success"].(bool) {
			url := "/dashboard/user/" + userId
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		}
	}
}
