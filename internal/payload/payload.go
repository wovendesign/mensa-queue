package payload

import (
	"fmt"
	parsers "mensa-queue/internal"

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

type LocalRecipe struct {
	Locales []RecipesLocales
	Additive       *[]Additive
	Nutrients      *[]Nutrient
	Recipe 	   Recipe
}

type Recipe struct {
	ID             uint `gorm:"primaryKey"`
	Diet           Diet
	PriceStudents  *float64
	PriceEmployees *float64
	PriceGuests    *float64
	MensaProvider  int64       `gorm:"column:mensa_provider_id"`
}

type RecipesRel struct {
	ID 	 uint `gorm:"primaryKey"`
	ParentID uint `gorm:"column:parent_id"`
	Path string
	AdditivesID *uint `gorm:"column:additives_id"`
	AllergensID *uint `gorm:"column:allergens_id"`
}

type RecipesLocales struct {
	ID uint `gorm:"primaryKey"`
	Title string `gorm:"column:title"`
	Locale string `gorm:"column:_locale"`
	RecipeID uint `gorm:"column:_parent_id"`
}

type Additive struct {
	ID   uint `gorm:"primaryKey"`
}

type AdditiveLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	AdditiveID uint `gorm:"column:_parent_id"`
}

type Allergen struct {
	ID   uint `gorm:"primaryKey"`
}

type AllergenLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	AllergenID uint `gorm:"column:_parent_id"`
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
	UnitId         uint         `gorm:"column:unit_id"`
	Unit           NutrientUnit `gorm:"foreignKey:unit_id"`
	Recommendation *string
}

type NutrientLabelsLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
}

type NutrientUnit struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

func InsertRecipe(recipe LocalRecipe) {
	// Database connection
	dsn := "host=127.0.0.1 user=mensauser password=postgres dbname=mensahhub port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if len(recipe.Locales) == 0 {
		fmt.Println("No locales provided")
		return
	}

	// Check if recipe already exists
	// (Title and MensaProvider are unique together)
	// If it does not exist, insert it
	nutrients := recipe.Nutrients

	var count RecipesLocales
	db.FirstOrInit(&count, recipe.Locales[0])

	if count.ID == 0 {
		// Create Recipe without title
		if err := db.Create(&recipe.Recipe).Error; err != nil {
			fmt.Println("Error inserting recipe:", err)
			panic(err)
		}
	} else {
		if err := db.Where(Recipe{
			ID: count.RecipeID,
			MensaProvider: 1,
		}).Assign(recipe.Recipe).FirstOrCreate(&recipe.Recipe).Error; err != nil {
			fmt.Println("Error inserting recipe:", err)
			panic(err)
		}
	}

	for _ , locale := range recipe.Locales {
		locale.RecipeID = recipe.Recipe.ID
		if err := db.FirstOrCreate(&locale, locale).Error; err != nil {
			fmt.Println("Error inserting locale:", err)
			panic(err)
		}
	}

	for _, nutrient := range *nutrients {
		nutrient.RecipeID = recipe.Recipe.ID
		_, err := insertNutrient(nutrient, db)
		if err != nil {
			fmt.Println("Error inserting nutrient:", err)
			panic(err)
		}
	}

	fmt.Println("Recipe inserted successfully", recipe.Recipe.ID)
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

func insertAllergen(allergen Allergen, name parsers.LocalizedString, recipe Recipe, db *gorm.DB) (*Allergen, error) {

	localAllergenDE := AllergenLocale{
		Name: *name.ValueDE,
	}
	var nameDE AllergenLocale
	db.FirstOrInit(&nameDE, localAllergenDE)
	localAllergenEN := AllergenLocale{
		Name: *name.ValueEN,
	}
	var nameEN AllergenLocale
	db.FirstOrInit(&nameEN, localAllergenEN)

	if nameDE.ID == 0 {
		// Create Recipe without title
		if err := db.Create(&allergen).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}

		// Create Locales
		nameDE.AllergenID = allergen.ID
		if err := db.Create(nameDE).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}
		nameEN.AllergenID = allergen.ID
		if err := db.Create(nameEN).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}
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

func insertAdditive(additive Additive, name parsers.LocalizedString, recipe Recipe, db *gorm.DB) (*Additive, error) {
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
