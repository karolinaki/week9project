package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/karolinaki/week9project/blob/main/routes.go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func main() {
	router := routes.SetupRoutes()

	log.Fatal(http.ListenAndServe(":8080", router))
	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"), // Replace with your desired region
		Credentials: credentials.NewStaticCredentials("ASIAW7OMQ3DYAVT7CCEQ", "PO1Y82NGpTGWs9wqmffpwIwT03MU3JnYYbikjMzE", "IQoJb3JpZ2luX2VjEDYaCXVzLWVhc3QtMSJFMEMCH2UNCg+Kyh1zpvC95d/noeyIZ+DaEUFJ6diIrRksuWICICpIDFKiXzFkcZKg1Y3MO7PoZ0e9caW6qeINaBE2U8+AKvUCCF8QABoMNDc5ODU0NjQ3NTM2Igw4dixdg6jALZu9veQq0gKqvLQyg9mt5L6iZ/PLN9FftjXeGBDzpR9D/riZnSOr3//SoRBifyjFTYgR/jzi4F9fTNkcwtumgolNDRgfKd/7+5H4NE6cw6v60LB+FjUpavK2D4nW0E2kfnaNmmdDQK3QKFV1Xc5HJhHf7+Rf7v23kZpZXPm9o7h0jcso0ewysr5qsQGsndCNxl/+DELkThSvhIQrihyZALB7iRQpefzJfC/dNTCVdpFQnqeR5NDWBsctmvAL2CZQ2AnyFPdaDcS+EMHoNU1T51AfCKgVC5NvVupeO7G8GvEVVLhwGopwSrlns9e7BdyaFeMUuBUogCpsgsmgUxB4PKfv6yeoCdX770XKQj/dh+A8irRUzoTZiyVOjMX+C5U0DfralpjWEBL2hRMBYx61C7q2cICdwGo9t2hdMybraVA6kJp+WNXOIZs4wVyIVyYdBXZNDcONN7XNFTDO0ZijBjqpAbU0O96vtqNoj5qH5ppTpnr9zO5m9fhqeZTA78W81pEFcWbSar6LeAJNVaqeiUYFBI3XEBoB8Kpn6v82gLfg9RCsEy2g7vQlKvJ8+nnE4L4t0PxTTxrQdG8o0l/Ami77uuOyYfBsFA6Nyamhk6GTyu/VD0SP/cVOiiRCu3KQK4R3O1mfxLXfw6mze21eJx0L6IU0eM2oKX7RclnOc9DXw3rwNLtXNJt3KaY="),
	})

	if err != nil {
		// Handle session creation error
		fmt.Println("Error creating session:", err)
		return
	}
	svc := dynamodb.New(sess)
	seedData, err := readSeedData("seed.json")
	if err != nil {
		log.Fatalf("Failed to read seed data: %v", err)
	}

	// Insert seed data into DynamoDB
	err = insertSeedData(svc, seedData)
	if err != nil {
		log.Fatalf("Failed to insert seed data into DynamoDB: %v", err)
	}

	fmt.Println("Seed data inserted successfully")

	fmt.Println("testing our functions!")
	getRecipeByID(svc, "2")
	createRecipe(svc, Recipe{ID: "5", Name: "bananas", Ingredients: []string{"bananas"}, CookingTime: 1, Image: "na"})
	deleteRecipeByID(svc, "5")

}

func readSeedData(filename string) ([]Recipe, error) {
	file, err := os.Open("seed.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open seed file: %v", err)
	}
	defer file.Close()

	var seedData []Recipe
	err = json.NewDecoder(file).Decode(&seedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode seed data: %v", err)
	}

	return seedData, nil
}

// Insert seed data into DynamoDB
func insertSeedData(svc dynamodbiface.DynamoDBAPI, seedData []Recipe) error {
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

func getRecipeByID(svc *dynamodb.DynamoDB, recipeID string) (*Recipe, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("YourTableName"), // Replace with your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(recipeID), // Replace with the ID of the recipe you want to retrieve
			},
		},
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil // Recipe not found
	}

	recipe := new(Recipe)
	err = dynamodbattribute.UnmarshalMap(result.Item, recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func createRecipe(svc *dynamodb.DynamoDB, recipe Recipe) error {
	av, err := dynamodbattribute.MarshalMap(recipe)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("YourTableName"), // Replace with your DynamoDB table name
		Item:      av,
	}

	_, err = svc.PutItem(input)
	return err
}

func deleteRecipeByID(svc *dynamodb.DynamoDB, recipeID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("YourTableName"), // Replace with your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(recipeID), // Replace with the ID of the recipe you want to delete
			},
		},
	}

	_, err := svc.DeleteItem(input)
	return err
}
