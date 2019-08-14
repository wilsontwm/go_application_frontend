package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/csrf"
	"log"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"app_frontend/controllers"
	"app_frontend/middleware"
)

func main() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	router := mux.NewRouter()
	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("csrf_token")),
		// To be removed in production in https
		csrf.Secure(false),
	)
	router.Use(csrfMiddleware)

	// Routes
	nonAuthenticatedRoutes := router.PathPrefix("").Subrouter()
	nonAuthenticatedRoutes.Use(middleware.Logging(), middleware.IsLoggedIn())
	
	// Pages routes
	nonAuthenticatedRoutes.HandleFunc("/", controllers.WelcomePage).Methods("GET").Name("welcome")
	nonAuthenticatedRoutes.HandleFunc("/noaccess", controllers.Custom403Page).Name("error_403")

	// Login / register routes	
	nonAuthenticatedRoutes.HandleFunc("/login", controllers.LoginPage).Methods("GET").Name("login")
	nonAuthenticatedRoutes.HandleFunc("/login", controllers.LoginSubmit).Methods("POST").Name("login_submit")
	nonAuthenticatedRoutes.HandleFunc("/logout", controllers.LogoutSubmit).Methods("POST").Name("logout")
	nonAuthenticatedRoutes.HandleFunc("/signup", controllers.SignupPage).Methods("GET").Name("signup")
	nonAuthenticatedRoutes.HandleFunc("/signup", controllers.SignupSubmit).Methods("POST").Name("signup_submit")
	nonAuthenticatedRoutes.HandleFunc("/resendactivation", controllers.ResendActivationPage).Methods("GET").Name("resend_activation")	
	nonAuthenticatedRoutes.HandleFunc("/resendactivation", controllers.ResendActivationSubmit).Methods("POST").Name("resend_activation_submit")
	nonAuthenticatedRoutes.HandleFunc("/activate/{code}", controllers.ActivateAccountPage).Methods("GET").Name("activate_account")	
	nonAuthenticatedRoutes.HandleFunc("/forgetpassword", controllers.ForgetPasswordPage).Methods("GET").Name("forget_password")	
	nonAuthenticatedRoutes.HandleFunc("/forgetpassword", controllers.ForgetPasswordSubmit).Methods("POST").Name("forget_password_submit")
	nonAuthenticatedRoutes.HandleFunc("/resetpassword/{code}", controllers.ResetPasswordPage).Methods("GET").Name("reset_password")	
	nonAuthenticatedRoutes.HandleFunc("/resetpassword/{code}", controllers.ResetPasswordSubmit).Methods("POST").Name("reset_password_submit")
	
	authenticatedRoutes := router.PathPrefix("/dashboard").Subrouter()
	authenticatedRoutes.Use(middleware.Logging(), middleware.CheckAuth())	
	authenticatedRoutes.HandleFunc("", controllers.DashboardPage).Methods("GET").Name("dashboard")

	// Profile routes
	profileRoutes := authenticatedRoutes.PathPrefix("/profile").Subrouter()
	profileRoutes.HandleFunc("/edit", controllers.EditProfilePage).Methods("GET").Name("profile_edit")	
	profileRoutes.HandleFunc("/edit", controllers.EditProfileSubmit).Methods("POST").Name("profile_edit_submit")
	profileRoutes.HandleFunc("/edit/password", controllers.EditPasswordSubmit).Methods("POST").Name("profile_edit_password_submit")
	profileRoutes.HandleFunc("/upload/picture", controllers.UploadPictureSubmit).Methods("POST").Name("profile_upload_picture_submit")
	profileRoutes.HandleFunc("/delete/picture", controllers.DeletePictureSubmit).Methods("POST").Name("profile_delete_picture_submit")
	
	// Company routes
	companyRoutes := authenticatedRoutes.PathPrefix("/company").Subrouter()
	companyRoutes.HandleFunc("", controllers.CompanyIndexPage).Methods("GET").Name("company_index")
	companyRoutes.HandleFunc("/store", controllers.CompanyCreateSubmit).Methods("POST").Name("company_store")
	companyRoutes.HandleFunc("/getUniqueSlug", controllers.CompanyGetUniqueSlugJson).Methods("GET").Name("company_get_unique_slug_json")
	companyRoutes.HandleFunc("/{id}/show", controllers.CompanyShowPage).Methods("GET").Name("company_show")
	companyRoutes.HandleFunc("/{id}/show/json", controllers.CompanyShowJson).Methods("GET").Name("company_show_json")
	companyRoutes.HandleFunc("/{id}/update", controllers.CompanyEditSubmit).Methods("POST").Name("company_edit_submit")
	companyRoutes.HandleFunc("/{id}/delete", controllers.CompanyDeleteSubmit).Methods("POST").Name("company_delete_submit")
	companyRoutes.HandleFunc("/{id}/invite", controllers.CompanyInviteSubmit).Methods("POST").Name("company_invite_submit")

	// Asset files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	router.PathPrefix("/storage/").Handler(http.StripPrefix("/storage/", http.FileServer(http.Dir("./storage/"))))
	
	// Custom 404 page
	router.NotFoundHandler = http.HandlerFunc(controllers.Custom404Page)

	port := os.Getenv("port")
	if port == "" {
		port = "9000"
	}

	log.Println("Frontend Server started and running at port", port)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Fatal(http.ListenAndServe(":" + port, handlers.CORS(headers, methods, origins)(router)))
}
