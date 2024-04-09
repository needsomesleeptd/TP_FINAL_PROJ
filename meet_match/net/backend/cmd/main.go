package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	auth_handler "test_backend_frontend/internal/http-server/handlers/auth"
	"test_backend_frontend/internal/http-server/handlers/cards"
	scroll2 "test_backend_frontend/internal/http-server/handlers/scroll"
	sessions_handler "test_backend_frontend/internal/http-server/handlers/session"
	"test_backend_frontend/internal/middleware/auth_middleware"
	"test_backend_frontend/internal/model"
	"test_backend_frontend/internal/models/models_da"
	auth_service "test_backend_frontend/internal/services/auth"
	repo_adapter "test_backend_frontend/internal/services/auth/user_repo/user_repo_ad"
	"test_backend_frontend/internal/services/scroll"
	postgres2 "test_backend_frontend/internal/services/scroll/scroll_repo/postgres"
	sessions "test_backend_frontend/internal/sessions"
	"test_backend_frontend/pkg/auth_utils"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	CONN_POSTGRES_STR = "host=localhost user=any1 password=1 database=meetmatch_db port=5432" //TODO:: export through parameters
	POSTGRES_CFG      = postgres.Config{DSN: CONN_POSTGRES_STR}
	MODEL_ROUTE       = "http://0.0.0.0:5000/rec"
	SESSION_PATH      = "localhost:6379"
)

func main() {
	model, err := model.New(MODEL_ROUTE)
	if err != nil {
		fmt.Println("Error with model")
		os.Exit(1)
	}
	var sessionManager *sessions.SessionManager
	sessionManager, err = sessions.NewSessionManager(SESSION_PATH, "", 0)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db, err := gorm.Open(postgres.New(POSTGRES_CFG), &gorm.Config{})
	db.AutoMigrate(models_da.User{}) //TODO:: this is a hack, fix this
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// TODO : add config

	//auth service
	userRepo := repo_adapter.NewUserRepositoryAdapter(db)
	hasher := auth_utils.NewPasswordHashCrypto()
	tokenHandler := auth_utils.NewJWTTokenHandler()
	userService := auth_service.NewAuthService(userRepo, hasher, tokenHandler, auth_service.SECRET)
	router := chi.NewRouter()

	// Scroll service
	scrollRepo := postgres2.NewScrollRepository(db)
	scrollManager := scroll.NewScrollUseCase(scrollRepo, sessionManager)

	authMiddleware := (func(h http.Handler) http.Handler {
		return auth_middleware.JwtAuthMiddleware(h, auth_service.SECRET, tokenHandler)
	})
	router.Group(func(r chi.Router) { //group for which auth middleware is required
		r.Use(authMiddleware)
		r.Get("/cards", cards.New(model))
		r.Post("/sessions", sessions_handler.SessionCreatePage(sessionManager))
		r.Post("/sessions/{id}", sessions_handler.SessionGetData(sessionManager))
		r.Patch("/sessions/{id}", sessions_handler.SessionAdduser(sessionManager))
		r.Put("/sessions/{id}", sessions_handler.SessionModifyuser(sessionManager))
		r.Get("/sessions/{id}/check_match", scroll2.NewCheckHandler(scrollManager))
		r.Post("/sessions/{id}/scroll", scroll2.NewScrollFactRegistrateHandler(scrollManager, tokenHandler))
	})

	//auth
	router.Post("/user/SignUp", auth_handler.SignUp(userService))
	router.Post("/user/SignIn", auth_handler.SignIn(userService))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
		IdleTimeout:  40 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("error with server")
		}
	}()

	<-done
}
