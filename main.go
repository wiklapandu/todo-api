package main

import (
	"log"

	"github.com/wiklapandu/todo-api/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	routes.Routes()
}
