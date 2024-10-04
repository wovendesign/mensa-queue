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


type LocalizedString struct {
	ValueDE string
	ValueEN string
}

type LocalizedValue[T any] struct {
	DE T
	EN T
}

type LocalNutrient struct {
	Nutrient Nutrient
	Name LocalizedString
}

type LocalRecipe struct {
	Locales []RecipesLocales
	Allergen       *[][]AllergensLocale
	Additives      *[][]AdditivesLocale
	Nutrients      *[]LocalNutrient
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

type AdditivesLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	AdditiveID uint `gorm:"column:_parent_id"`
}

type Allergen struct {
	ID   uint `gorm:"primaryKey"`
}

type AllergensLocale struct {
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
	NutrientLabelID uint `gorm:"column:_parent_id"`
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

	for _, nutrient := range *recipe.Nutrients {
		nutrient.Nutrient.RecipeID = recipe.Recipe.ID
		_, err := insertNutrient(nutrient.Nutrient, nutrient.Name, db)
		if err != nil {
			fmt.Println("Error inserting nutrient:", err)
			panic(err)
		}
	}

	if (recipe.Allergen != nil) {
		for _, allergens := range *recipe.Allergen {
			_, err := insertAllergen(allergens, recipe.Recipe, db)
			if err != nil {
				fmt.Println("Error inserting allergen:", err)
				panic(err)
			}
		}
	}

	if (recipe.Additives != nil) {
		for _, additives := range *recipe.Additives {
			_, err := insertAdditive(additives, recipe.Recipe, db)
			if err != nil {
				fmt.Println("Error inserting additive:", err)
				panic(err)
			}
		}
	}

	fmt.Println("Recipe inserted successfully", recipe.Recipe.ID)
}

func insertNutrient(nutrient Nutrient, name LocalizedString, db *gorm.DB) (*Nutrient, error) {
	var unit NutrientUnit
	var nutrientValue NutrientValue
	var nutrientLabel NutrientLabel

	// Insert Nutrient Unit
	if err := db.FirstOrCreate(&unit, NutrientUnit{Name: nutrient.NutrientLabel.Unit.Name}).Error; err != nil {
		fmt.Println("Error inserting unit:", err)
		return nil, err
	}
	nutrient.NutrientLabel.Unit = unit
	nutrient.NutrientLabel.UnitId = unit.ID

	// Insert Nutrient Value
	if err := db.FirstOrCreate(&nutrientValue, NutrientValue{Value: nutrient.NutrientValue.Value}).Error; err != nil {
		fmt.Println("195")
		return nil, err
	}
	nutrient.NutrientValue = nutrientValue
	nutrient.NutrientValueID = nutrientValue.ID

	// Insert Nutrient Label
	localNutrientDE := NutrientLabelsLocale{
		Name: name.ValueDE,
		Locale: "de",
	}
	var nameDE NutrientLabelsLocale
	db.FirstOrInit(&nameDE, localNutrientDE)
	localNutrientEN := NutrientLabelsLocale{
		Name: name.ValueEN,
		Locale: "en",
	}
	var nameEN NutrientLabelsLocale
	db.FirstOrInit(&nameEN, localNutrientEN)


	if nameDE.ID == 0 {
		nutrientLabel = nutrient.NutrientLabel
		if err := db.Create(&nutrientLabel).Error; err != nil {
			fmt.Println("Couldnt Create Nutrient Label")
			return nil, err
		}

		// Create Locales
		nameDE.NutrientLabelID = nutrientLabel.ID
		fmt.Println("NameDE:", nameDE)
		if err := db.Create(&nameDE).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}
		nameEN.NutrientLabelID = nutrientLabel.ID
		if err := db.Create(&nameEN).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}
	} else {
		if err := db.Where(NutrientLabel{
			ID: nameDE.NutrientLabelID,
		}).First(&nutrientLabel).Error; err != nil {
			return nil, err
		}
	}

	nutrient.NutrientLabel = nutrientLabel
	nutrient.NutrientLabelID = nutrientLabel.ID

	// Insert Nutrient
	if err := db.FirstOrCreate(&nutrient, nutrient).Error; err != nil {
		return nil, err
	}

	return &nutrient, nil
}

func insertAllergen(allergens []AllergensLocale, recipe Recipe, db *gorm.DB) (*Allergen, error) {
	localAllergenDE := allergens[0]
	var nameDE AllergensLocale
	db.FirstOrInit(&nameDE, localAllergenDE)
	localAllergenEN := allergens[1]
	var nameEN AllergensLocale
	db.FirstOrInit(&nameEN, localAllergenEN)

	var allergen Allergen

	if nameDE.ID == 0 {
		// Create Recipe without title
		if err := db.Create(&allergen).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}

		// Create Locales
		nameDE.AllergenID = allergen.ID
		if err := db.Create(&nameDE).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}
		nameEN.AllergenID = allergen.ID
		if err := db.Create(&nameEN).Error; err != nil {
			fmt.Println("Error inserting allergen:", err)
			return nil, err
		}
	} else {
		if err := db.Where(Allergen{
			ID: nameDE.AllergenID,
		}).First(&allergen).Error; err != nil {
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

func insertAdditive(additive []AdditivesLocale, recipe Recipe, db *gorm.DB) (*Additive, error) {
	localAdditiveDE := additive[0]
	var nameDE AdditivesLocale
	db.FirstOrInit(&nameDE, localAdditiveDE)
	localAdditiveEN := additive[1]
	var nameEN AdditivesLocale
	db.FirstOrInit(&nameEN, localAdditiveEN)

	var additiveModel Additive

	if nameDE.ID == 0 {
		// Create Recipe without title
		if err := db.Create(&additiveModel).Error; err != nil {
			fmt.Println("Error inserting additive:", err)
			return nil, err
		}

		// Create Locales
		nameDE.AdditiveID = additiveModel.ID
		if err := db.Create(&nameDE).Error; err != nil {
			fmt.Println("Error inserting additive:", err)
			return nil, err
		}
		nameEN.AdditiveID = additiveModel.ID
		if err := db.Create(&nameEN).Error; err != nil {
			fmt.Println("Error inserting additive:", err)
			return nil, err
		}
	} else {
		if err := db.Where(Additive{
			ID: nameDE.AdditiveID,
		}).First(&additiveModel).Error; err != nil {
			return nil, err
		}
	}

	// Attach the additive to the recipe using the association method
	rel := RecipesRel{
		ParentID: recipe.ID,
		Path: "additives",
		AllergensID: nil,
		AdditivesID: &additiveModel.ID,
	}
	if err := db.Create(&rel).Error; err != nil {
		fmt.Println("Error inserting recipe additive:", err)
		return nil, err
	}

	return &additiveModel, nil
}
