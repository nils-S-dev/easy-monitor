package main

import (
	"easy-monitor/internal/api"
	"easy-monitor/internal/schedule"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing with system environment variables")
	}
	schedule.Init()
	http.ListenAndServe(":8080", api.NewRouter())
}
