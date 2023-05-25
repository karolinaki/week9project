package models

type Recipe struct {
	ID          string   `json:"ID"`
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	CookingTime int      `json:"cookingTime"`
	Image       string   `json:"imgURL"`
	// Add more attributes as needed
}
