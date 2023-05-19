package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Define the route handlers
func GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
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
