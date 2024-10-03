package payload

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Diet is a custom type representing the diet categories
type Diet string

// Define constants for each diet type
const (
	DietVegan      Diet = "vegan"
	DietVegetarian Diet = "vegetarian"
	DietMeat       Diet = "meat"
	DietFish       Diet = "fish"
)

type Recipe struct {
	ID             uint `gorm:"primaryKey"`
	Title          string
	Diet           Diet
	PriceStudents  *float64
	PriceEmployees *float64
	PriceGuests    *float64
	Nutrients      *[]Nutrient
	MensaProvider  int64       `gorm:"column:mensa_provider_id"`
	Additive       *[]Additive `gorm:"many2many:recipe_rels;"`
	AllergensID       *[]Allergen `gorm:"many2many:recipe_rels;"`
}

type RecipesRel struct {
	ID 	 uint `gorm:"primaryKey"`
	ParentID uint `gorm:"column:parent_id"`
	Path string
	AdditivesID *uint `gorm:"column:additives_id"`
	AllergensID *uint `gorm:"column:allergens_id"`
}

type Additive struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Allergen struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Nutrient struct {
	ID              uint          `gorm:"primaryKey"`
	NutrientValueID uint          `gorm:"column:nutrient_value_id"`
	NutrientValue   NutrientValue `gorm:"foreignKey:nutrient_value_id"`
	NutrientLabelID uint          `gorm:"column:nutrient_label_id"`
	NutrientLabel   NutrientLabel `gorm:"foreignKey:nutrient_label_id"`
	RecipeID        uint          `gorm:"column:recipe_id"`
	Recipe          Recipe        `gorm:"foreignKey:recipe_id"`
}

type NutrientValue struct {
	ID    uint    `gorm:"primaryKey"`
	Value float64 `gorm:"unique"`
}

type NutrientLabel struct {
	ID             uint         `gorm:"primaryKey"`
	Name           string       `gorm:"unique"`
	UnitId         uint         `gorm:"column:unit_id"`
	Unit           NutrientUnit `gorm:"foreignKey:unit_id"`
	Recommendation *string
}

type NutrientUnit struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

func InsertRecipe(recipe Recipe) {
	// Database connection
	dsn := "host=127.0.0.1 user=mensauser password=postgres dbname=mensahhub port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Check if recipe already exists
	// (Title and MensaProvider are unique together)
	// If it does not exist, insert it
	nutrients := recipe.Nutrients
	recipe.Nutrients = nil



	if err := db.Where(Recipe{
		Title: recipe.Title,
		MensaProvider: 1,
	}).Assign(Recipe{
		Diet: recipe.Diet,
		PriceStudents: recipe.PriceStudents,
		PriceEmployees: recipe.PriceEmployees,
		PriceGuests: recipe.PriceGuests,
	}).FirstOrCreate(&recipe).Error; err != nil {
		fmt.Println("Error inserting recipe:", err)
		panic(err)
	}

	for _, nutrient := range *nutrients {
		nutrient.RecipeID = recipe.ID
		_, err := insertNutrient(nutrient, db)
		if err != nil {
			fmt.Println("Error inserting nutrient:", err)
			panic(err)
		}
	}

	fmt.Println("Recipe inserted successfully", recipe.ID)
}

func insertNutrient(nutrient Nutrient, db *gorm.DB) (*Nutrient, error) {
	var unit NutrientUnit
	var nutrientValue NutrientValue
	var nutrientLabel NutrientLabel

	// Insert Nutrient Unit
	if err := db.FirstOrCreate(&unit, NutrientUnit{Name: nutrient.NutrientLabel.Unit.Name}).Error; err != nil {
		return nil, err
	}
	nutrient.NutrientLabel.Unit = unit
	nutrient.NutrientLabel.UnitId = unit.ID

	// Insert Nutrient Value
	if err := db.FirstOrCreate(&nutrientValue, NutrientValue{Value: nutrient.NutrientValue.Value}).Error; err != nil {
		return nil, err
	}
	nutrient.NutrientValue = nutrientValue
	nutrient.NutrientValueID = nutrientValue.ID

	// Insert Nutrient Label
	if err := db.FirstOrCreate(&nutrientLabel, nutrient.NutrientLabel).Error; err != nil {
		return nil, err
	}
	nutrient.NutrientLabel = nutrientLabel
	nutrient.NutrientLabelID = nutrientLabel.ID

	// Insert Nutrient
	if err := db.FirstOrCreate(&nutrient, nutrient).Error; err != nil {
		return nil, err
	}

	return &nutrient, nil
}

func insertAllergen(allergen Allergen, recipe Recipe, db *gorm.DB) (*Allergen, error) {
	if err := db.FirstOrCreate(&allergen, allergen).Error; err != nil {
		fmt.Println("Error inserting allergen:", err)
		return nil, err
	}

	// Attach the allergens to the recipe using the association method
	rel := RecipesRel{
		ParentID: recipe.ID,
		Path: "allergens",
		AllergensID: &allergen.ID,
		AdditivesID: nil,
	}
	if err := db.Create(&rel).Error; err != nil {
		fmt.Println("Error inserting recipe allergen:", err)
		return nil, err
	}

	return &allergen, nil
}

func insertAdditive(additive Additive, recipe Recipe, db *gorm.DB) (*Additive, error) {
	if err := db.FirstOrCreate(&additive, additive).Error; err != nil {
		fmt.Println("Error inserting additive:", err)
		return nil, err
	}

	// Attach the additives to the recipe using the association method
	rel := RecipesRel{
		ParentID: recipe.ID,
		Path: "additives",
		AllergensID: nil,
		AdditivesID: &additive.ID,
	}
	if err := db.Create(&rel).Error; err != nil {
		fmt.Println("Error inserting recipe additive:", err)
		return nil, err
	}

	return &additive, nil
}
