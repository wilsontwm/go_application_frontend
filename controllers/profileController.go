package controllers

import (
	util "app_frontend/utils"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Display the edit profile page for the current user
var EditProfilePage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	// Set the URL path
	restURL.Path = "/api/dashboard/profile/get"
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

		data := map[string]interface{}{
			"title":          "Edit Profile",
			"appName":        appName,
			"appVersion":     appVersion,
			"name":           name,
			"picture":        picture,
			"year":           year,
			"user":           resp["data"].(map[string]interface{}),
			"countries":      resp["countries"].([]interface{}),
			"genders":        resp["genders"].([]interface{}),
			csrf.TemplateTag: csrf.TemplateField(r),
		}

		data, err = util.InitializePage(w, r, store, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = templates.ExecuteTemplate(w, "edit_profile_html", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Post request for editing the profile for the current user
var EditProfileSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/dashboard/profile/edit"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	phone := strings.TrimSpace(r.Form.Get("phone"))
	city := strings.TrimSpace(r.Form.Get("city"))
	country, _ := strconv.ParseInt(r.Form.Get("country"), 10, 32)
	gender, _ := strconv.ParseInt(r.Form.Get("gender"), 10, 32)
	bio := strings.TrimSpace(r.Form.Get("bio"))
	dateFormat := "02 Jan 2006" // dd M yyyy
	birthday, _ := time.Parse(dateFormat, r.Form.Get("birthday"))

	// Set the input data
	jsonData := map[string]interface{}{
		"name":     name,
		"phone":    phone,
		"city":     city,
		"country":  country,
		"gender":   gender,
		"birthday": birthday,
		"bio":      bio,
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

		if resp["success"].(bool) {
			// Need to reset the cookie that store name
			userData := resp["data"].(map[string]interface{})
			SetCookieHandler(w, r, "name", userData["name"].(string))
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// Post request to edit the password for the current user
var EditPasswordSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/dashboard/profile/edit/password"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	password := strings.TrimSpace(r.Form.Get("password"))
	retype_password := strings.TrimSpace(r.Form.Get("retype_password"))

	// Check if the retype password matches
	if password != retype_password {
		session.AddFlash("Retype password does not match.", "errors")
		session.Save(r, w)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

		return
	}

	// Set the input data
	jsonData := map[string]interface{}{
		"password": password,
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

// Post request to upload profile picture for the current user
var UploadPictureSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// Set the URL path
	restURL.Path = "/api/dashboard/profile/upload/picture"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")
	picture := ReadCookieHandler(w, r, "picture")

	//Parse the multipart form, 3 << 20 specifies a maximum of 3 MB files
	r.ParseMultipartForm(3 << 20)

	// Get the file
	file, handler, err := r.FormFile("picture")

	if err != nil || handler.Size > 3000000 {
		session.AddFlash("Error retrieving the file. Please make sure that the file size does not exceed 3MB.", "errors")
		session.Save(r, w)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

		return
	}

	defer file.Close()

	extension := filepath.Ext(handler.Filename)
	name := handler.Filename[0:len(handler.Filename)-len(extension)] + time.Now().String()
	hash := md5.New()
	hash.Write([]byte(fmt.Sprint(name)))
	finalFileName := hex.EncodeToString(hash.Sum(nil)) + "*" + extension

	// Create a tempory file within the folder "storage/profile"
	tempFile, err := ioutil.TempFile("storage/profile", finalFileName)
	filePath := tempFile.Name()

	// Convert the string to slash
	filePath = strings.Replace(filePath, "\\", "/", -1)

	if err != nil {
		session.AddFlash("Error uploading the file.", "errors")
		session.Save(r, w)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

		return
	}

	defer tempFile.Close()

	// Send file path to the database
	// Set the input data
	jsonData := map[string]interface{}{
		"profilePicture": "/" + filePath,
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

		if resp["success"].(bool) {

			// Read all the contents into byte array
			fileBytes, _ := ioutil.ReadAll(file)

			// Write the byte array to temporary file
			tempFile.Write(fileBytes)

			// Remove the first character of the picture as it consists of slash
			_, i := utf8.DecodeRuneInString(picture)
			picture = picture[i:]

			if !strings.HasPrefix(picture, "assets") {
				// Remove the existing file
				go os.Remove(picture)
			}

			// Need to reset the cookie that store name
			userData := resp["data"].(map[string]interface{})

			profilePicture := defaultProfilePic // default profile picture
			if userData["profilePicture"] != nil && userData["profilePicture"] != "" {
				profilePicture = userData["profilePicture"].(string)
			}

			SetCookieHandler(w, r, "picture", profilePicture)
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// Request to delete the profile picture
var DeletePictureSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// Set the URL path
	restURL.Path = "/api/dashboard/profile/delete/picture"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")
	picture := ReadCookieHandler(w, r, "picture")

	// Set the input data
	jsonData := make(map[string]interface{})

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

		if resp["success"].(bool) {
			// Remove the first character of the picture as it consists of slash
			_, i := utf8.DecodeRuneInString(picture)
			picture = picture[i:]
			if !strings.HasPrefix(picture, "assets") {
				// Remove the existing file
				go os.Remove(picture)
			}

			// Need to reset the cookie that store name
			userData := resp["data"].(map[string]interface{})

			profilePicture := defaultProfilePic // default profile picture
			if userData["profilePicture"] != nil && userData["profilePicture"] != "" {
				profilePicture = userData["profilePicture"].(string)
			}

			SetCookieHandler(w, r, "picture", profilePicture)
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// Display the profile page of user
var ProfileShowPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()

	session, err := util.GetSession(store, w, r)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	userId := vars["id"]

	// Set the URL path
	restURL.Path = "/api/dashboard/user/" + userId
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

		if resp["success"].(bool) {
			var profile map[string]interface{}

			if _, ok := resp["data"]; ok {
				profile = resp["data"].(map[string]interface{})

				if profile["profilePicture"] == nil || profile["profilePicture"] == "" {
					profile["profilePicture"] = defaultProfilePic // default profile picture
				}

				// Set default country name as blank first
				profile["countryName"] = ""
				if _, ok := resp["countries"]; ok {
					countries := resp["countries"].([]interface{})
					for _, c := range countries {
						country := c.(map[string]interface{})
						if country["CountryId"] == profile["country"] {
							profile["countryName"] = country["Name"]
							break
						}
					}
				}

				// Set default gender as blank first
				profile["genderType"] = ""
				if _, ok := resp["genders"]; ok {
					genders := resp["genders"].([]interface{})
					for _, g := range genders {
						gender := g.(map[string]interface{})
						if gender["GenderId"] == profile["gender"] {
							profile["genderType"] = gender["Sex"]
							break
						}
					}
				}
			}

			data := map[string]interface{}{
				"title":          "Profile",
				"appName":        appName,
				"appVersion":     appVersion,
				"name":           name,
				"picture":        picture,
				"year":           year,
				"profile":        profile,
				"isPersonal":     (userId == ReadCookieHandler(w, r, "id")),
				csrf.TemplateTag: csrf.TemplateField(r),
			}

			data, err = util.InitializePage(w, r, store, data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = templates.ExecuteTemplate(w, "profile_show_html", data)

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
