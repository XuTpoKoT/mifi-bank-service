package router

import (
	"bank-service/internal/middleware"
	"database/sql"

	"bank-service/internal/handler"
	"bank-service/internal/repository"
	"bank-service/internal/service"

	"github.com/gorilla/mux"
)

func Setup(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Logging)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := handler.NewAccountHandler(accountService)

	cardRepo := repository.NewCardRepository(db)
	cardService := service.NewCardService(
		cardRepo,
		accountRepo,
	)
	cardHandler := handler.NewCardHandler(
		cardService,
	)

	creditRepo := repository.NewCreditRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	creditService := service.NewCreditService(
		db,
		creditRepo,
		accountRepo,
		scheduleRepo,
	)
	creditHandler := handler.NewCreditHandler(
		creditService,
	)

	r.HandleFunc(
		"/register",
		authHandler.Register,
	).Methods("POST")

	r.HandleFunc(
		"/login",
		authHandler.Login,
	).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.Auth)

	protected.HandleFunc(
		"/accounts",
		accountHandler.Create,
	).Methods("POST")

	protected.HandleFunc(
		"/accounts/topup",
		accountHandler.TopUp,
	).Methods("POST")

	protected.HandleFunc(
		"/transfer",
		accountHandler.Transfer,
	).Methods("POST")

	protected.HandleFunc(
		"/cards",
		cardHandler.Create,
	).Methods("POST")

	protected.HandleFunc(
		"/accounts",
		accountHandler.GetAll,
	).Methods("GET")

	protected.HandleFunc(
		"/cards",
		cardHandler.GetAll,
	).Methods("GET")

	protected.HandleFunc(
		"/credits",
		creditHandler.Create,
	).Methods("POST")

	protected.HandleFunc(
		"/credits/{creditId}/schedule",
		creditHandler.GetSchedule,
	).Methods("GET")

	return r
}
