package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"github.com/karolinaki/week9project/db"
	"github.com/karolinaki/week9project/models"
)

func getRecipeByID(svc *dynamodb.DynamoDB, recipeID string) (*models.Recipe, error) {
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

	recipe := new(models.Recipe)
	err = dynamodbattribute.UnmarshalMap(result.Item, recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func createRecipe(svc *dynamodb.DynamoDB, recipe models.Recipe) error {
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

// Define the route handlers
func GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the recipe ID from the request parameters
	vars := mux.Vars(r)
	recipeID := vars["id"]

	// Call the getRecipeByID function to retrieve the recipe from DynamoDB
	recipe, err := getRecipeByID(svc*dynamodb.DynamoDB, recipeID)
	if err != nil {
		log.Printf("Failed to get recipe: %v", err)
		http.Error(w, "Failed to get recipe", http.StatusInternalServerError)
		return
	}

	// Write the recipe JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request, svc *dynamodb.DynamoDB) {
	// Extract the recipe ID from the request parameters
	vars := mux.Vars(r)
	recipeID := vars["id"]

	// Call the deleteRecipeByID function to delete the recipe from DynamoDB
	err := deleteRecipeByID(svc, recipeID)
	if err != nil {
		log.Printf("Failed to delete recipe: %v", err)
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	// Write the success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Recipe deleted successfully"))
}

// Define other route handlers for POST and UPDATE operations similarly

// Define the router and register the route handlers
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Define the routes and associate them with the corresponding route handlers
	r.HandleFunc("/recipes/{id}", GetRecipeHandler).Methods("GET")
	r.HandleFunc("/recipes/{id}", DeleteRecipeHandler).Methods("DELETE")
	// Add routes for POST and UPDATE operations

	return r
}
