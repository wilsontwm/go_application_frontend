package main

import (
	"app_frontend/controllers"
	"app_frontend/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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

	// Invitation routes (incoming)
	invitedRoutes := authenticatedRoutes.PathPrefix("/invite/incoming").Subrouter()
	invitedRoutes.HandleFunc("", controllers.IndexInvitationFromCompany).Methods("GET")
	invitedRoutes.HandleFunc("/{id}/respond", controllers.RespondCompanyInvitationRequestSubmit).Methods("POST")

	// Company routes
	companyRoutes := authenticatedRoutes.PathPrefix("/company").Subrouter()
	companyRoutes.HandleFunc("", controllers.CompanyIndexPage).Methods("GET").Name("company_index")
	companyRoutes.HandleFunc("/store", controllers.CompanyCreateSubmit).Methods("POST").Name("company_store")
	companyRoutes.HandleFunc("/getUniqueSlug", controllers.CompanyGetUniqueSlugJson).Methods("GET").Name("company_get_unique_slug_json")
	companyRoutes.HandleFunc("/users/search/all", controllers.CompanyUsersSearchPage).Methods("GET").Name("company_users_search_page")
	companyRoutes.HandleFunc("/users/search", controllers.CompanyUsersSearchJson).Methods("GET").Name("company_users_search_json")
	companyRoutes.HandleFunc("/{id}/show", controllers.CompanyShowPage).Methods("GET").Name("company_show")
	companyRoutes.HandleFunc("/{id}/show/json", controllers.CompanyShowJson).Methods("GET").Name("company_show_json")
	companyRoutes.HandleFunc("/{id}/update", controllers.CompanyEditSubmit).Methods("POST").Name("company_edit_submit")
	companyRoutes.HandleFunc("/{id}/delete", controllers.CompanyDeleteSubmit).Methods("POST").Name("company_delete_submit")
	companyRoutes.HandleFunc("/{id}/users", controllers.CompanyUsersListJson).Methods("GET").Name("company_users_list_json")
	companyRoutes.HandleFunc("/{id}/visit", controllers.CompanyVisitSubmit).Methods("POST").Name("company_visit_submit")

	// Company invitation request routes
	companyRoutes.HandleFunc("/{id}/invite", controllers.CompanyInviteSubmit).Methods("POST").Name("company_invite_submit")
	companyRoutes.HandleFunc("/{id}/invite/list", controllers.CompanyInviteListJson).Methods("GET").Name("company_invite_list_json")
	companyRoutes.HandleFunc("/{id}/invite/multiple/resend", controllers.CompanyInvitationResendMultipleSubmit).Methods("POST").Name("company_invite_resend_multiple_submit")
	companyRoutes.HandleFunc("/{id}/invite/multiple/delete", controllers.CompanyInvitationDeleteMultipleSubmit).Methods("POST").Name("company_invite_delete_multiple_submit")
	companyRoutes.HandleFunc("/{id}/invite/{invitationID}/resend", controllers.CompanyInvitationResendSubmit).Methods("POST").Name("company_invite_resend_submit")
	companyRoutes.HandleFunc("/{id}/invite/{invitationID}/delete", controllers.CompanyInvitationDeleteSubmit).Methods("POST").Name("company_invite_delete_submit")

	// User routes
	userRoutes := authenticatedRoutes.PathPrefix("/user").Subrouter()
	userRoutes.HandleFunc("/{id}", controllers.ProfileShowPage).Methods("GET")

	// Post routes
	postRoutes := authenticatedRoutes.PathPrefix("/post").Subrouter()
	postRoutes.HandleFunc("/create", controllers.PostCreatePage).Methods("GET").Name("post_create_page")
	postRoutes.HandleFunc("/store", controllers.PostCreateSubmit).Methods("POST").Name("post_store")
	postRoutes.HandleFunc("/{id}/show", controllers.PostShowPage).Methods("GET").Name("post_show_page")
	postRoutes.HandleFunc("/{id}/edit", controllers.PostEditPage).Methods("GET").Name("post_edit_page")

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

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
