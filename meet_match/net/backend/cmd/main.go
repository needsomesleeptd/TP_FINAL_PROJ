package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	auth_handler "test_backend_frontend/internal/http-server/handlers/auth"
	"test_backend_frontend/internal/http-server/handlers/cards"
	feedback_handler "test_backend_frontend/internal/http-server/handlers/feedback"
	sessions_handler "test_backend_frontend/internal/http-server/handlers/session"
	"test_backend_frontend/internal/middleware/auth_middleware"
	"test_backend_frontend/internal/models/models_da"
	rec_model_client "test_backend_frontend/internal/rec-model-client"
	auth_service "test_backend_frontend/internal/services/auth"
	repo_adapter "test_backend_frontend/internal/services/auth/user_repo/user_repo_ad"
	postgres3 "test_backend_frontend/internal/services/cards/repository/postgres"
	feedback_service "test_backend_frontend/internal/services/feedback"
	"test_backend_frontend/internal/services/feedback/feedback_repo"
	match_repo_adap "test_backend_frontend/internal/services/match/matchRepo/matchRepoAd"
	"test_backend_frontend/internal/services/scroll"
	sessions "test_backend_frontend/internal/sessions"
	"test_backend_frontend/pkg/auth_utils"
	"time"

	scroll2 "test_backend_frontend/internal/http-server/handlers/scroll"
	postgres2 "test_backend_frontend/internal/services/scroll/scroll_repo/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	CONN_POSTGRES_STR = "host='postgres' user=any1 password='1' database='meetmatch_db' port=5432 sslmode=disable" //TODO:: export through parameters
	POSTGRES_CFG      = postgres.Config{DSN: CONN_POSTGRES_STR}
	MODEL_ROUTE       = "http://python-flask-app:5000/rec"
	SESSION_PATH      = "redis-db:6379"
)

func main() {
	db, err := gorm.Open(postgres.New(POSTGRES_CFG), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		//os.Exit(1)
		time.Sleep(5 * time.Second)
	}

	db, err = gorm.Open(postgres.New(POSTGRES_CFG), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(5 * time.Second)
	}
	cardRepo := postgres3.NewCardRepo(db)
	model, err := rec_model_client.New(MODEL_ROUTE, cardRepo)
	if err != nil {
		fmt.Println("Error with model")
		os.Exit(1)
	}
	tokenHandler := auth_utils.NewJWTTokenHandler()
	var sessionManager *sessions.SessionManager
	sessionManager, err = sessions.NewSessionManager(SESSION_PATH, "", 0, tokenHandler, auth_service.SECRET)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db.AutoMigrate(models_da.User{}) //TODO: this is a hack, fix this
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// TODO: add config

	//auth service
	userRepo := repo_adapter.NewUserRepositoryAdapter(db)
	hasher := auth_utils.NewPasswordHashCrypto()

	userService := auth_service.NewAuthService(userRepo, hasher, tokenHandler, auth_service.SECRET)
	router := chi.NewRouter()

	matchRepo := match_repo_adap.NewFeedbackRepo(db)
	//	Scroll service
	scrollRepo := postgres2.NewScrollRepository(db)
	scrollManager := scroll.NewScrollUseCase(scrollRepo, sessionManager, cardRepo, matchRepo)

	//	feedback service
	feedbackRepo := feedback_repo.NewFeedbackRepo(db)
	feedbackService := feedback_service.NewFeedbackService(feedbackRepo)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(cors.Handler)

	authMiddleware := (func(h http.Handler) http.Handler {
		return auth_middleware.JwtAuthMiddleware(h, auth_service.SECRET, tokenHandler)
	})
	router.Group(func(r chi.Router) { //group for which auth middleware is required
		r.Use(authMiddleware)
		r.Post("/cards", cards.New(model, tokenHandler))
		r.Post("/sessions", sessions_handler.SessionCreatePage(sessionManager))
		r.Post("/sessions/{id}", sessions_handler.SessionsGetSessionData(sessionManager))
		r.Patch("/sessions/{id}", sessions_handler.SessionAdduser(sessionManager))
		r.Put("/sessions/{id}", sessions_handler.SessionModifyuser(sessionManager))
		r.Post("/sessions/getUser", sessions_handler.SessionGetUserSessions(sessionManager))
		r.Post("/sessions/getSessionUsers", sessions_handler.SessionGetData(sessionManager))
		r.Delete("/sessions/{id}", sessions_handler.SessionDeleteUser(sessionManager))

		r.Post("/sessions/{id}/check_match", scroll2.NewCheckHandler(scrollManager))
		r.Post("/sessions/{id}/scroll", scroll2.NewScrollFactRegistrateHandler(scrollManager, tokenHandler, cardRepo))

		r.Post("/feedback/has_gone", feedback_handler.SaveFeedback(feedbackService))

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
