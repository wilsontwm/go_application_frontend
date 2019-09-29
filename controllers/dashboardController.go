package controllers

import (
	"net/http"
	util "app_frontend/utils"
	"time"
	"github.com/gorilla/csrf"
)

var DashboardPage = func(w http.ResponseWriter, r *http.Request) {
	name := ReadCookieHandler(w, r, "name")
	picture := ReadCookieHandler(w, r, "picture")
	year := time.Now().Year()
	data := map[string]interface{}{
		"title": "Dashboard",
		"appName": appName,
		"appVersion": appVersion,
		"name": name,
		"picture": picture,
		"year": year,
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "dashboard_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
