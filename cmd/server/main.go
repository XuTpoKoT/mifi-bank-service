package main

import (
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

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	r := router.Setup(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server started on :" + port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
