package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/erazr/test-task/config"
	db "github.com/erazr/test-task/db"
	"github.com/erazr/test-task/models"
	services "github.com/erazr/test-task/services"

	_ "github.com/erazr/test-task/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	userService *services.UserService
	database    *db.MongoDB
}

func NewHandler(cfg config.Config, database *db.MongoDB, userService *services.UserService) *Handler {
	return &Handler{
		userService: userService,
		database:    database,
	}
}

// @title Test task BackDev
// @version 1.0
// @description This is a simple API project
// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

func (h *Handler) RunHttp(ctx context.Context, port string, swaggerPath string) {
	if port == "" {
		port = ":8080"
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle(swaggerPath, httpSwagger.WrapHandler)

	mux.HandleFunc("/api/v1/authenticate", h.Authenticate)
	mux.HandleFunc("/api/v1/register", h.Register)
	mux.HandleFunc("/api/v1/refresh", h.Refresh)

	log.Println("Server started on port", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()
}

// @Summary Generate pair of tokens
// @Description Generate pair of tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param guid query string true "User guid"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /authenticate [get]
func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	guid := r.URL.Query().Get("guid")
	if guid == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Guid is required")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := h.userService.Authenticate(r.Context(), guid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err, "Error while authenticating user")
		return
	}

	json.NewEncoder(w).Encode(response)
}

// @Summary Register new user
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User object to be added"
// @Success 201 {string} string "User created"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body models.User

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err, "Error while decoding request body")
		return
	}

	err = h.userService.Register(r.Context(), body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err, "Error while registering user")
		return
	}

	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("User created"))
}

// @Summary Refresh token
// @Description Refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body models.RefreshToken true "Refresh token"
// @Success 200 {object} models.Tokens
// @Failure 500 {object} string
// @Router /refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var token models.RefreshToken

	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err, "Error while decoding request body")
		return
	}

	tokens, err := h.userService.Refresh(token.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err, "Invalid token")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tokens)
}
