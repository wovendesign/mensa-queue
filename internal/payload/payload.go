package payload

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Mensa int
const (
	NeuesPalais Mensa = 9600
	Golm Mensa = 9601
	Teltow Mensa = 9602
	Griebnitzsee Mensa = 9603
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
	FeaturesID *uint `gorm:"column:features_id"`
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

func InsertRecipe(recipe LocalRecipe, date time.Time) {
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

	InsertServing(date, NeuesPalais, recipe.Recipe.ID, db)
}

type Serving struct {
	ID 	  uint `gorm:"primaryKey"`
	Date  time.Time
	MensaID uint `gorm:"column:mensa_id"`
	RecipeID uint `gorm:"column:recipe_id"`
}

func InsertServing(date time.Time, mensa Mensa, recipeID uint, db *gorm.DB) {
	fmt.Printf("Inserting serving for %v on %s for recipe%d\n", mensa, date, recipeID)
	mensaMap := map[Mensa]uint{
		NeuesPalais: 1,
	}

	serving := Serving{
		Date: date,
		MensaID: mensaMap[mensa],
		RecipeID: recipeID,
	}

	// Check if serving already exists
	var count Serving
	db.FirstOrInit(&count, serving)

	if count.ID == 0 {
		err := db.Create(&serving).Error
		if err != nil {
			fmt.Println("Error inserting serving:", err)
			panic(err)
		}
	}
	fmt.Printf("Serving inserted successfully: %d\n", serving.ID)
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
	allergen := Allergen{}

	allergenEntity, err := insertEntityWithLocales[Allergen, AllergensLocale](db, allergen, allergens, func(locale AllergensLocale) error {
		locale.AllergenID = allergen.ID
		return db.Create(&locale).Error
	})
	if err != nil {
		return nil, err
	}

	// Attach the allergens to the recipe using the association method
	rel := RecipesRel{
		ParentID: recipe.ID,
		Path: "allergens",
		AllergensID: &allergenEntity.ID,
		AdditivesID: nil,
	}
	if err := db.Create(&rel).Error; err != nil {
		fmt.Println("Error inserting recipe allergen:", err)
		return nil, err
	}

	return &allergen, nil
}

func insertAdditive(additiveLocales []AdditivesLocale, recipe Recipe, db *gorm.DB) (*Additive, error) {
	additive := Additive{}

	// Create or find the Additive entity along with its DE/EN locales
	additiveEntity, err := insertEntityWithLocales(db, additive, additiveLocales, func(locale AdditivesLocale) error {
		locale.AdditiveID = additive.ID
		return db.Create(&locale).Error
	})
	if err != nil {
		return nil, err
	}

	// Attach the additive to the recipe using the association method
	rel := RecipesRel{
		ParentID: recipe.ID,
		Path: "additives",
		AllergensID: nil,
		AdditivesID: &additiveEntity.ID,
	}
	if err := db.Create(&rel).Error; err != nil {
		fmt.Println("Error inserting recipe additive:", err)
		return nil, err
	}

	return additiveEntity, nil
}

type Feature struct {
	ID uint `gorm:"primaryKey"`
	MensaProviderID uint `gorm:"column:mensa_provider_id"`
}

type FeatureLocale struct {
	ID uint `gorm:"primaryKey"`
	FeatureID uint `gorm:"column:_parent_id"`
	Locale string `gorm:"column:_locale"`
	Name string `gorm:"column:name"`
}


func insertFeature(featureLocales []FeatureLocale, recipe Recipe, db *gorm.DB) (*Feature, error) {
	feature := Feature{
		MensaProviderID: 1,
	}

	// Create or find the Feature entity along with its DE/EN locales
	featureEntity, err := insertEntityWithLocales(db, feature, featureLocales, func(locale FeatureLocale) error {
		locale.FeatureID = feature.ID
		return db.Create(&locale).Error
	})
	if err != nil {
		return nil, err
	}

	// Attach the feature to the recipe using the association method
	rel := RecipesRel{
		ParentID: recipe.ID,
		Path: "features",
		AllergensID: nil,
		AdditivesID: nil,
		FeaturesID: &featureEntity.ID,
	}
	if err := db.Create(&rel).Error; err != nil {
		fmt.Println("Error inserting recipe feature:", err)
		return nil, err
	}

	return featureEntity, nil
}

// findOrCreateLocale abstracts the repeated logic of checking and creating locales.
func findOrCreateLocale[T any](db *gorm.DB, locale T, createEntity func() (T, error)) (T, error) {
	var existingLocale T

	// Check if the locale already exists in the database
	if err := db.FirstOrInit(&existingLocale, locale).Error; err != nil {
		fmt.Println("Error finding locale:", err)
		return existingLocale, err
	}

	// If the locale doesn't exist (ID is zero), call the entity creation function
	// createEntity is a function that creates the entity if it doesn't exist
	if getID(existingLocale) == 0 {
		createdLocale, err := createEntity()
		if err != nil {
			fmt.Println("Error creating entity:", err)
			return createdLocale, err
		}
		return createdLocale, nil
	}

	return existingLocale, nil
}

// Helper function to extract ID (assuming all structs have an ID field)
func getID[T any](v T) uint {
	switch x := any(v).(type) {
	case AdditivesLocale:
		return x.ID
	case AllergensLocale:
		return x.ID
	case NutrientLabelsLocale:
		return x.ID
	// Add more cases as needed for other types
	default:
		return 0
	}
}

func insertLocales[T any](db *gorm.DB, locales []T, localeCreateFn func(T) error) error {
	for _, locale := range locales {
		_, err := findOrCreateLocale(db, locale, func() (T, error) {
			if err := localeCreateFn(locale); err != nil {
				return locale, err
			}
			return locale, nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func insertEntityWithLocales[T any, L any](db *gorm.DB, entity T, locales []L, localeCreateFn func(L) error) (*T, error) {
	// Create the entity if it doesn't exist
	entity, err := findOrCreateLocale(db, entity, func() (T, error) {
		if err := db.Create(&entity).Error; err != nil {
			return entity, err
		}
		return entity, nil
	})
	if err != nil {
		return nil, err
	}

	// Insert associated locales
	if err := insertLocales(db, locales, localeCreateFn); err != nil {
		return nil, err
	}

	return &entity, nil
}
