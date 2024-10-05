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

type Language int
const (
	DE Language = iota + 1
	EN
)

func (l Language) String() string {
	return [...]string{"", "de", "en"}[l]
}

type LocalizedString map[Language]string

type LocalizedValue[T any] struct {
	DE T
	EN T
}

type LocalNutrient struct {
	Nutrient Nutrient
	Name LocalizedString
}

type LocalRecipe struct {
	Locales []Locale
	Allergen       *[][]Locale
	Additives      *[][]Locale
	Features	   *[][]Locale
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

type RecipesLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
type AdditivesLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
type AllergensLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
type FeaturesLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
type NutrientsLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
type Locale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}


type Additive struct {
	ID   uint `gorm:"primaryKey"`
}

type Allergen struct {
	ID   uint `gorm:"primaryKey"`
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

type NutrientUnit struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

func InsertRecipe(recipe LocalRecipe, date time.Time, language []Language) {
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
	count := RecipesLocale{
		ID: recipe.Recipe.ID,
		Name: recipe.Locales[0].Name,
		Locale: recipe.Locales[0].Locale,
		ParentID: recipe.Recipe.ID,
	}
	db.FirstOrInit(&count, count)

	if count.ID == 0 {
		// Create Recipe without title
		if err := db.Create(&recipe.Recipe).Error; err != nil {
			fmt.Println("Error inserting recipe:", err)
			panic(err)
		}
	} else {
		if err := db.Where(Recipe{
			ID: count.ParentID,
			MensaProvider: 1,
		}).Assign(recipe.Recipe).FirstOrCreate(&recipe.Recipe).Error; err != nil {
			fmt.Println("Error inserting recipe:", err)
			panic(err)
		}
	}

	for _ , locale := range recipe.Locales {
		locale.ParentID = recipe.Recipe.ID
		if err := db.Table("recipes_locales").FirstOrCreate(&locale, locale).Error; err != nil {
			fmt.Println("Error inserting locale:", err)
			panic(err)
		}
	}

	for _, nutrient := range *recipe.Nutrients {
		nutrient.Nutrient.RecipeID = recipe.Recipe.ID
		_, err := insertNutrient(nutrient.Nutrient, nutrient.Name, db, language)
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

	if (recipe.Features != nil) {
		for _, feature := range *recipe.Features {
			_, err = insertFeature(feature, recipe.Recipe, db)
			if err != nil {
				fmt.Println("Error inserting feature: ", err)
				return
			}
		}
	}

	InsertServing(date, NeuesPalais, recipe.Recipe.ID, db)
}

type Serving struct {
	ID 	  uint `gorm:"primaryKey"`
	Date  time.Time
	MensaID uint `gorm:"column:mensa_id"`
	RecipeID uint `gorm:"column:recipe_id"`
}

func InsertServing(date time.Time, mensa Mensa, recipeID uint, db *gorm.DB) {
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
}

type NutrientLabelsLocale struct {
	ID             uint `gorm:"primaryKey"`
	Name 		 string
	ParentID       uint `gorm:"column:_parent_id"`
	Locale 	   string `gorm:"column:_locale"`
}

func insertNutrient(nutrient Nutrient, name LocalizedString, db *gorm.DB, languages []Language) (*Nutrient, error) {
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
	nutrientLocales := make([]NutrientLabelsLocale, 0)
	for _, language := range languages {
		localNutrient := NutrientLabelsLocale{
			Name: name[language],
			Locale: language.String(),
		}
		var name NutrientLabelsLocale
		db.FirstOrInit(&name, localNutrient)
		nutrientLocales = append(nutrientLocales, name)
	}

	if &nutrientLocales[0] == nil {
		fmt.Println("No Nutrient Locales generated :(")
		return nil, fmt.Errorf("No nutrient locales generated")
	}

	if nutrientLocales[0].ID == 0 {
		nutrientLabel = nutrient.NutrientLabel
		if err := db.Create(&nutrientLabel).Error; err != nil {
			fmt.Println("Couldnt Create Nutrient Label")
			return nil, err
		}

		// Create Locales
		for _, language := range nutrientLocales {
			language.ParentID = nutrientLabel.ID
			if err := db.Create(&language).Error; err != nil {
				fmt.Println("Error inserting allergen:", err)
				return nil, err
			}
		}
	} else {
		if err := db.Where(NutrientLabel{
			ID: nutrientLocales[0].ParentID,
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

func insertAllergen(allergens []Locale, recipe Recipe, db *gorm.DB) (*Allergen, error) {
	allergen := Allergen{}

	allergenEntity, err := insertEntityWithLocales[Allergen](db, allergen, &allergens)
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

func insertAdditive(additiveLocales []Locale, recipe Recipe, db *gorm.DB) (*Additive, error) {
	additive := Additive{}

	// Create or find the Additive entity along with its DE/EN locales
	additiveEntity, err := insertEntityWithLocales(db, additive, &additiveLocales)
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

func insertFeature(featureLocales []Locale, recipe Recipe, db *gorm.DB) (*Feature, error) {
	feature := Feature{
		MensaProviderID: 1,
	}

	// Create or find the Feature entity along with its DE/EN locales
	featureEntity, err := insertEntityWithLocales(db, feature, &featureLocales)
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

func getID[T any](entity T) uint {
	// return entity.ID
	switch x := any(entity).(type) {
		case Additive:
			return x.ID
		case Allergen:
			return x.ID
		case Feature:
			return x.ID
		case Nutrient:
			return x.ID
		default:
			return 0
	}
}

// findOrCreateLocale abstracts the repeated logic of checking and creating locales.
func findOrCreateLocale[T any](db *gorm.DB, locale *Locale, entity T) error {
	var existingLocale Locale

	// Check if the locale already exists in the database
	// if err := db.FirstOrInit(&existingLocale, locale).Error; err != nil {
	// 	fmt.Println("Error finding locale:", err)
	// 	return err
	// }

	switch any(entity).(type) {
		case Additive:
			_locale := AdditivesLocale{
				ID: locale.ID,
				Name: locale.Name,
				Locale: locale.Locale,
				ParentID: locale.ParentID,
			}
			// Check if the locale already exists in the database
			if err := db.FirstOrCreate(&_locale, _locale).Error; err != nil {
				fmt.Println("Error finding locale:", err)
				return err
			}
			locale.ID = _locale.ID
			existingLocale = Locale(_locale)
		case Allergen:
			_locale := AllergensLocale{
				ID: locale.ID,
				Name: locale.Name,
				Locale: locale.Locale,
				ParentID: locale.ParentID,
			}
			// Check if the locale already exists in the database
			if err := db.FirstOrCreate(&_locale, _locale).Error; err != nil {
				fmt.Println("Error finding locale:", err)
				return err
			}
			fmt.Println("Locale ID:", _locale.ID)
			locale.ID = _locale.ID
			existingLocale = Locale(_locale)
		case Feature:
			_locale := FeaturesLocale{
				ID: locale.ID,
				Name: locale.Name,
				Locale: locale.Locale,
				ParentID: locale.ParentID,
			}
			// Check if the locale already exists in the database
			if err := db.FirstOrCreate(&_locale, _locale).Error; err != nil {
				fmt.Println("Error finding locale:", err)
				return err
			}
			locale.ID = _locale.ID
			existingLocale = Locale(_locale)
		case Nutrient:
			_locale := NutrientsLocale{
				ID: locale.ID,
				Name: locale.Name,
				Locale: locale.Locale,
				ParentID: locale.ParentID,
			}
			// Check if the locale already exists in the database
			if err := db.FirstOrCreate(&_locale, _locale).Error; err != nil {
				fmt.Println("Error finding locale:", err)
				return err
			}
			locale.ID = _locale.ID
			existingLocale = Locale(_locale)
		default:
			return fmt.Errorf("Table name not found")
	}

	// If the locale doesn't exist (ID is zero), call the entity creation function
	// createEntity is a function that creates the entity if it doesn't exist
	if existingLocale.ID == 0 {
		// No Locales found
		// Create entity, then create locale
		if err := db.Create(&entity).Error; err != nil {
			return err
		}
		existingLocale.ParentID = getID(entity)

		return nil
	}

	return nil
}

func insertEntityWithLocales[T any](db *gorm.DB, entity T, locales *[]Locale) (*T, error) {
	for _, locale := range *locales {
		err := findOrCreateLocale(db, &locale, entity)
		if err != nil {
			return nil, err
		}
		// locales[i] = *locale
	}
	return &entity, nil
}
