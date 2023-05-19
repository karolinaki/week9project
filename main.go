package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/karolinaki/week9project/db"
	"github.com/karolinaki/week9project/models"
	"github.com/karolinaki/week9project/routes"
)

func main() {
	// Read seed data
	seedData, err := readSeedData("seed.json")
	if err != nil {
		log.Fatalf("Failed to read seed data: %v", err)
	}

	// Initialize DynamoDB client
	sess := db.NewSession()
	svc := dynamodb.New(sess)

	// Insert seed data into DynamoDB
	err = insertSeedData(svc, seedData)
	if err != nil {
		log.Fatalf("Failed to insert seed data into DynamoDB: %v", err)
	}

	// Set up routes
	r := routes.NewRouter(svc)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

func readSeedData(filename string) ([]models.Recipe, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open seed file: %v", err)
	}
	defer file.Close()

	var seedData []models.Recipe
	err = json.NewDecoder(file).Decode(&seedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode seed data: %v", err)
	}

	return seedData, nil
}

// Insert seed data into DynamoDB
func insertSeedData(svc dynamodbiface.DynamoDBAPI, seedData []models.Recipe) error {
	for _, recipe := range seedData {
		av, err := dynamodbattribute.MarshalMap(recipe)
		if err != nil {
			return fmt.Errorf("failed to marshal recipe data: %v", err)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("YourTableName"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			return fmt.Errorf("failed to insert item into DynamoDB: %v", err)
		}
	}

	return nil
}
