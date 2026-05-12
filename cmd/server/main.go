package main

import (
	"bank-service/internal/logger"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"bank-service/internal/db"
	"bank-service/internal/router"
)

func main() {

	// загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using system env")
	}
	logger.Init()

	database, err := db.Connect()
	if err != nil {
		logger.Log.Fatal(err)
	}

	r := router.Setup(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.Info("server started on :" + port)

	logger.Log.Fatal(http.ListenAndServe(":"+port, r))
}
