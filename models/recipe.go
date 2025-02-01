package models

type Recipe struct {
	ID          string    // Unique ID for the recipe
	Name        string    // Name of the recipe
	Description string    // Description of the recipe
	Ingredients []string  // List of ingredients
	Servings    []Serving // List of servings for different canteens
}

type Serving struct {
	CanteenID string  // ID of the canteen
	Portion   string  // Portion size (e.g., "Small", "Large")
	Quantity  float64 // Quantity of ingredients for this serving
}
