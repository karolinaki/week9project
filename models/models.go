package models

type Recipe struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	CookingTime int      `json:"cookingTime"`
	Image       string   `json:"imgURL"`
	// Add more attributes as needed
}
