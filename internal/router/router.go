package router

import (
	"database/sql"

	"bank-service/internal/handler"
	"bank-service/internal/repository"
	"bank-service/internal/service"

	"github.com/gorilla/mux"
)

func Setup(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r.HandleFunc(
		"/register",
		authHandler.Register,
	).Methods("POST")

	r.HandleFunc(
		"/login",
		authHandler.Login,
	).Methods("POST")

	return r
}
