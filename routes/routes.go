package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"github.com/karolinaki/week9project/models"
)

// GetRecipeHandler handles the GET request to retrieve a recipe by ID.
func GetRecipeHandler(w http.ResponseWriter, r *http.Request, svc *dynamodb.DynamoDB) {
	// Extract the recipe ID from the request parameters
	vars := mux.Vars(r)
	recipeID := vars["id"]

	// Call the getRecipeByID function to retrieve the recipe from DynamoDB
	recipe, err := getRecipeByID(svc, recipeID)
	if err != nil {
		log.Printf("Failed to get recipe: %v", err)
		http.Error(w, "Failed to get recipe", http.StatusInternalServerError)
		return
	}

	// Write the recipe JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Define other route handler functions as needed

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

// Create a new router and set up routes
func NewRouter(svc *dynamodb.DynamoDB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetRecipeHandler(w, r, svc)
	}).Methods("GET")
	// Define other routes as needed

	return r
}
