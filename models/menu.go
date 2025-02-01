package models

type MenuItem struct {
	ID          string  // Unique ID for the menu item
	Name        string  // Name of the menu item
	Description string  // Description of the menu item
	Price       float64 // Price of the menu item
	RecipeID    string  // Reference to the recipe
}
