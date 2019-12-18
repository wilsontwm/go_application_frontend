package utils

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
	"reflect"
)

// Build json message
func Message(success bool, status float64, message string, errors []string) map[string]interface{} {
	return map[string]interface{}{"success": success, "status": status, "message": message, "errors": errors}
}

// Return json response
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_, hasData := data["status"]
	if hasData {
		w.WriteHeader(int(data["status"].(float64)))
	}
	json.NewEncoder(w).Encode(data)
}

// Send a post request to the url
func SendPostRequest(url string, data map[string]interface{}) (response *http.Response, err error) {
	requestBody, err := json.Marshal(data)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err = client.Do(request)

	return
}

// Send a request to the url with authentication
func SendAuthenticatedRequest(url string, requestType string, authToken string, data map[string]interface{}) (response *http.Response, err error) {
	requestBody, err := json.Marshal(data)

	request, _ := http.NewRequest(requestType, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}

	response, err = client.Do(request)

	return
}

// Initialize a page
func InitializePage(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, data map[string]interface{}) (output map[string]interface{}, err error) {
	session, err := GetSession(store, w, r)
	errorMessages := session.Flashes("errors")
	successMessage := session.Flashes("success")
	session.Save(r, w)

	flash := map[string]interface{}{
		"errors":  errorMessages,
		"success": successMessage,
	}
	output = MergeMapString(data, flash)

	// Get the user ID and check on the companies and selected company from redis
	id := ReadCookieHandler(w, r, "id")
	if id != "" {
		var companies []interface{}
		selectedCompany := GetActiveCompany(w, r, id)
		companiesRedis, err := RedisLRange("user:" + id + ";companies:")
		if err != nil {
			log.Println("Error reading companies from Redis:", err.Error())
		} else {
			for _, companyRedis := range companiesRedis {
				var comp interface{}
				json.Unmarshal(companyRedis, &comp)
				companies = append(companies, comp)
			}
		}

		companiesMap := map[string]interface{}{
			"dropdownCompanies":       companies,
			"dropdownSelectedCompany": selectedCompany,
		}

		output = MergeMapString(output, companiesMap)
	}

	return
}

// Get a session
func GetSession(store *sessions.CookieStore, w http.ResponseWriter, r *http.Request) (session *sessions.Session, err error) {
	err = godotenv.Load() //Load .env file
	sessionName := os.Getenv("session_name")
	session, err = store.Get(r, sessionName)
	return
}

// Build the error message
func GetErrorMessages(errors *[]string, err error) {
	for _, errz := range err.(validator.ValidationErrors) {
		// Build the custom errors here
		switch tag := errz.ActualTag(); tag {
		case "required":
			*errors = append(*errors, errz.StructField()+" is required.")
		case "email":
			*errors = append(*errors, errz.StructField()+" is an invalid email address.")
		case "min":
			if errz.Type().Kind() == reflect.String {
				*errors = append(*errors, errz.StructField()+" must be more than or equal to "+errz.Param()+" character(s).")
			} else {
				*errors = append(*errors, errz.StructField()+" must be larger than "+errz.Param()+".")
			}
		case "max":
			if errz.Type().Kind() == reflect.String {
				*errors = append(*errors, errz.StructField()+" must be lesser than or equal to "+errz.Param()+" character(s).")
			} else {
				*errors = append(*errors, errz.StructField()+" must be smaller than "+errz.Param()+".")
			}
		default:
			*errors = append(*errors, errz.StructField()+" is invalid.")
		}
	}

	return
}

// Merge two map string interface
func MergeMapString(mp1 map[string]interface{}, mp2 map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})
	for k, v := range mp1 {
		if _, ok := mp1[k]; ok {
			result[k] = v
		}
	}

	for k, v := range mp2 {
		if _, ok := mp2[k]; ok {
			result[k] = v
		}
	}

	return result
}

// Read cookie
func ReadCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string) (cookieValue string) {
	cookie, err := r.Cookie(cookieName)

	if err == nil {
		cookieValue = cookie.Value
	}

	return
}

// Get the current active company
func GetActiveCompany(w http.ResponseWriter, r *http.Request, userId string) interface{} {
	var selectedCompany interface{}
	selectedCompanyRedis, err := RedisGet("user:" + userId + ";selectedcompany:")
	if err != nil {
		log.Println("Error reading selected company from Redis:", err.Error())
	} else {
		err = json.Unmarshal(selectedCompanyRedis, &selectedCompany)
	}

	return selectedCompany
}

// Get the ID of the current active company
func GetActiveCompanyID(w http.ResponseWriter, r *http.Request, userId string) string {
	company := GetActiveCompany(w, r, userId)
	companyId := ""

	if companyJson, ok := company.(map[string]interface{}); ok {
		companyId = companyJson["ID"].(string)
	}

	return companyId
}
