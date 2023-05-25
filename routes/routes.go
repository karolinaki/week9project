package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	recipeID := vars["ID"]

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
func getRecipesAll(svc *dynamodb.DynamoDB, userProfile string) ([]*models.Recipe, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Recipes"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	recipes := []*models.Recipe{}
	for _, item := range result.Items {
		recipe := &models.Recipe{}
		err = dynamodbattribute.UnmarshalMap(item, recipe)
		if err != nil {
			return nil, err
		}

		// Check if the recipe ID is "No4" and the userProfile is not allowed to see it
		skipRecipe := recipe.ID == "4" && userProfile != "MrKrabbs" && userProfile != "SpongeBob"
		if skipRecipe {
			continue // Skip this recipe
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func createNewRecipe(svc *dynamodb.DynamoDB, recipe models.Recipe) error {
	// Generate a new unique ID for the recipe
	newID, err := generateNewID(svc)
	if err != nil {
		return err
	}
	recipe.ID = newID

	// Convert the recipe struct to a DynamoDB attribute value map
	av, err := dynamodbattribute.MarshalMap(recipe)
	if err != nil {
		return err
	}

	// Create the input for the PutItem operation
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Recipes"),
		Item:      av,
	}

	// Insert the recipe into DynamoDB
	_, err = svc.PutItem(input)
	return err
}
func generateNewID(svc *dynamodb.DynamoDB) (string, error) {
	// Get the count of existing recipes in the table
	input := &dynamodb.ScanInput{
		TableName:      aws.String("Recipes"),
		Select:         aws.String("COUNT"),
		ConsistentRead: aws.Bool(true),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return "", err
	}

	existingCount := *result.Count

	// Generate the new ID by adding 1 to the existing count
	newID := strconv.Itoa(int(existingCount) + 1)
	return newID, nil
}

func CreateRecipeHandler(w http.ResponseWriter, r *http.Request, svc *dynamodb.DynamoDB) {
	// Parse the request body to get the recipe data
	var recipe models.Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the createNewRecipe function to create the recipe in DynamoDB
	err = createNewRecipe(svc, recipe)
	if err != nil {
		log.Printf("Failed to create recipe: %v", err)
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}

	// Write a success response
	w.WriteHeader(http.StatusCreated)
}

func getRecipeByID(svc *dynamodb.DynamoDB, recipeID string) (*models.Recipe, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Recipes"), // Replace with your DynamoDB table name
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

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request, svc *dynamodb.DynamoDB) {
	// Extract the recipe ID from the request parameters
	vars := mux.Vars(r)
	recipeID := vars["id"]

	// Call the deleteRecipeByID function passing the recipeID
	err := deleteRecipeByID(svc, recipeID)
	if err != nil {
		// Handle the error
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	// Write a success response
	w.WriteHeader(http.StatusNoContent)
}

func deleteRecipeByID(svc *dynamodb.DynamoDB, recipeID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Recipes"), // Replace with your DynamoDB table name
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
	r.HandleFunc("/recipes/{ID}", func(w http.ResponseWriter, r *http.Request) {
		GetRecipeHandler(w, r, svc)
	}).Methods("GET")
	r.HandleFunc("/recipes/{ID}", func(w http.ResponseWriter, r *http.Request) { DeleteRecipeHandler(w, r, svc) }).Methods("DELETE")
	r.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) { CreateRecipeHandler(w, r, svc) }).Methods("POST")
	return r
}
