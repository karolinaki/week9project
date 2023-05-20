package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/karolinaki/week9project/models"
)

// NewDynamoDB initializes a new DynamoDB client and inserts seed data.
func NewDynamoDB(seedData []models.Recipe) *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	if err != nil {
		// Handle session creation error
		fmt.Println("Error creating session:", err)
	}

	svc := dynamodb.New(sess)

	// Insert seed data into DynamoDB
	err = insertSeedData(svc, seedData)
	if err != nil {
		log.Fatalf("Failed to insert seed data into DynamoDB: %v", err)
	}

	return svc
}

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

// Read seed data from file
func ReadSeedData(filename string) ([]models.Recipe, error) {
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
