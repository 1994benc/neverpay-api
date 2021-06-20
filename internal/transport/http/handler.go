package http

import (
	"1994benc/neverpay-user-service/internal/user"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"1994benc/neverpay-user-service/internal/transport/http/middleware"

	"github.com/gorilla/mux"
)

type Handler struct {
	Router      *mux.Router
	UserService *user.DefaultUserService
}

// Creates a new instance of Handler
func NewHandler(userService *user.DefaultUserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}

// Setup all routes
func (handler *Handler) SetupRoutes() {
	log.Println("Setting up routes")
	handler.Router = mux.NewRouter()

	// Middlewares
	handler.Router.Use(middleware.LoggingMiddleware)

	// All routes
	handler.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(GenericResponseModel{Message: "I am alive!"}); err != nil {
			panic(err)
		}
	})
	// User routes
	handler.Router.HandleFunc("/api/users", handler.GetUsers).Methods(http.MethodGet)
	handler.Router.HandleFunc("/api/users/signup", handler.SignUp).Methods(http.MethodPost)
	handler.Router.HandleFunc("/api/users/signin", handler.SignIn).Methods(http.MethodPost)
	// Only accessible from a server; requires a secret key to be passed in
	// Example: /api/verify?token=sometoken&scope=scope1,scope2
	handler.Router.HandleFunc("/api/verify", handler.VerifyToken).Methods(http.MethodGet)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var u user.UserModel
	err := u.FromJSON(r.Body)
	if err != nil {
		log.Error("Error parsing the JSON body: %s", err)
		http.Error(w, "Error parsing the JSON body", http.StatusBadRequest)
		return
	}
	userAlreadyExists := h.checkIfUserExists(u)
	if userAlreadyExists {
		http.Error(w, "Email already in use!", http.StatusBadRequest)
		return
	}

	err = u.HashPassword()
	if err != nil {
		http.Error(w, "Error generating hashed password", http.StatusInternalServerError)
		return
	}

	newUser, err := h.UserService.CreateUser(u)
	if err != nil {
		http.Error(w, "Error creating user! "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = newUser.ToJSON(w)
	if err != nil {
		http.Error(w, "Error writing to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var authDetails user.CredentialsModel
	error := authDetails.FromJSON(r.Body)
	if error != nil {
		http.Error(w, "Error parsing inputs", http.StatusBadRequest)
		return
	}

	u, err := h.UserService.FindUserByEmail(authDetails.Email)
	if u.Email == "" || err != nil {
		http.Error(w, "User not found in our system!", http.StatusForbidden)
		return
	}

	passwordsMatched := u.CheckPasswordHash(authDetails.Password)
	if !passwordsMatched {
		http.Error(w, "Password entered is incorrect!", http.StatusForbidden)
		return
	}

	validToken, err := h.UserService.GenerateJWT(u.Email, "basic")
	if err != nil {
		http.Error(w, "Error generating access token!", http.StatusInternalServerError)
		return
	}

	var token user.TokenModel
	token.Email = u.Email
	token.Role = u.Role
	token.TokenString = validToken
	err = token.ToJSON(w)
	if err != nil {
		http.Error(w, "Error parsing generated token!", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		http.Error(w, "Error getting users", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, "Token not provided", http.StatusUnauthorized)
		return
	}
	result := h.UserService.ValidateToken(token)
	result.ToJSON(w)
}

func (h *Handler) checkIfUserExists(userInstance user.UserModel) bool {
	u, err := h.UserService.FindUserByEmail(userInstance.Email)
	log.Printf("checkIfUserExists: %s", u.Email)
	return err == nil && u.Email != ""
}
