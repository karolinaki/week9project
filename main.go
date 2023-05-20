package main

import (
	"log"
	"net/http"

	"github.com/karolinaki/week9project/db"
	"github.com/karolinaki/week9project/routes"
)

func main() {
	// Read seed data
	seedData, err := db.ReadSeedData("seed.json")
	if err != nil {
		log.Fatalf("Failed to read seed data: %v", err)
	}

	// Initialize DynamoDB client
	svc := db.NewDynamoDB(seedData)

	// Set up routes
	r := routes.NewRouter(svc)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
