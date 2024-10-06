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
func (r Recipe) GetID() uint { return r.ID }
func (r Recipe) SetID(id uint) { r.ID = id }

type RecipesRel struct {
	ID 	 uint `gorm:"primaryKey"`
	ParentID uint `gorm:"column:parent_id"`
	Path string
	AdditivesID *uint `gorm:"column:additives_id"`
	AllergensID *uint `gorm:"column:allergens_id"`
	FeaturesID *uint `gorm:"column:features_id"`
}

type Locale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}

type EntityInterface interface {
	GetID() uint
	SetID(uint)
}

type LocaleInterface interface {
	GetID() uint
	GetName() string
	GetLocale() string
	GetParentID() uint
	SetParentID(uint)
}

type AdditivesLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
// Implement LocaleInterface methods for each locale type
func (l AdditivesLocale) GetID() uint      { return l.ID }
func (l AdditivesLocale) GetName() string  { return l.Name }
func (l AdditivesLocale) GetLocale() string { return l.Locale }
func (l AdditivesLocale) GetParentID() uint { return l.ParentID }
func (l *AdditivesLocale) SetParentID(id uint) { l.ParentID = id }

type RecipesLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
func (l RecipesLocale) GetID() uint      { return l.ID }
func (l RecipesLocale) GetName() string  { return l.Name }
func (l RecipesLocale) GetLocale() string { return l.Locale }
func (l RecipesLocale) GetParentID() uint { return l.ParentID }
func (l *RecipesLocale) SetParentID(id uint) { l.ParentID = id }

type AllergensLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
func (l AllergensLocale) GetID() uint      { return l.ID }
func (l AllergensLocale) GetName() string  { return l.Name }
func (l AllergensLocale) GetLocale() string { return l.Locale }
func (l AllergensLocale) GetParentID() uint { return l.ParentID }
func (l *AllergensLocale) SetParentID(id uint) { l.ParentID = id }

type FeaturesLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
func (l FeaturesLocale) GetID() uint      { return l.ID }
func (l FeaturesLocale) GetName() string  { return l.Name }
func (l FeaturesLocale) GetLocale() string { return l.Locale }
func (l FeaturesLocale) GetParentID() uint { return l.ParentID }
func (l *FeaturesLocale) SetParentID(id uint) { l.ParentID = id }

type NutrientsLocale struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Locale string `gorm:"column:_locale"`
	ParentID uint `gorm:"column:_parent_id"`
}
func (l NutrientsLocale) GetID() uint      { return l.ID }
func (l NutrientsLocale) GetName() string  { return l.Name }
func (l NutrientsLocale) GetLocale() string { return l.Locale }
func (l NutrientsLocale) GetParentID() uint { return l.ParentID }
func (l *NutrientsLocale) SetParentID(id uint) { l.ParentID = id }


type Additive struct {
	ID   uint `gorm:"primaryKey"`
}
func (a Additive) GetID() uint { return a.ID }
func (a Additive) SetID(id uint) { a.ID = id }

type Allergen struct {
	ID   uint `gorm:"primaryKey"`
}
func (a Allergen) GetID() uint { return a.ID }
func (a Allergen) SetID(id uint) { a.ID = id }


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
func (n NutrientLabel) GetID() uint { return n.ID }
func (n NutrientLabel) SetID(id uint) { n.ID = id }

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

	_, err := insertEntityWithLocales[Allergen](db, allergen, &allergens)
	if err != nil {
		return nil, err
	}

	// Attach the allergens to the recipe using the association method
	rel := RecipesRel{
		ParentID: recipe.ID,
		Path: "allergens",
		AllergensID: &allergens[0].ParentID,
		AdditivesID: nil,
	}
	if err := db.Where(rel).FirstOrCreate(&rel).Error; err != nil {
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
		AdditivesID: &additiveLocales[0].ParentID,
	}
	if err := db.Where(rel).FirstOrCreate(&rel).Error; err != nil {
		fmt.Println("Error inserting recipe additive:", err)
		return nil, err
	}

	return additiveEntity, nil
}

type Feature struct {
	ID uint `gorm:"primaryKey"`
	MensaProviderID uint `gorm:"column:mensa_provider_id"`
}
func (f Feature) GetID() uint { return f.ID }
func (f Feature) SetID(id uint) { f.ID = id }

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
		Path: "feature",
		AllergensID: nil,
		AdditivesID: nil,
		FeaturesID: &featureLocales[0].ParentID,
	}
	if err := db.Where(rel).FirstOrCreate(&rel).Error; err != nil {
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
func findOrCreateLocale[T any](db *gorm.DB, locale *Locale, entity T) (*Locale, error) {
	var existingLocale Locale
	var err error

	switch any(entity).(type) {
	case Additive:
		additiveLocale := AdditivesLocale{
			Name:     locale.Name,
			Locale:   locale.Locale,
			// ParentID: locale.ParentID,
		}
		err = db.Where(&additiveLocale).FirstOrInit(&additiveLocale).Error
		if err != nil {
			fmt.Println("Error finding/creating Additive locale:", err)
			return nil, err
		}
		existingLocale = Locale(additiveLocale)
	case Allergen:
		allergenLocale := AllergensLocale{
			Name:     locale.Name,
			Locale:   locale.Locale,
			// ParentID: locale.ParentID,
		}
		err = db.Where(&allergenLocale).FirstOrInit(&allergenLocale).Error
		if err != nil {
			fmt.Println("Error finding/creating Allergen locale:", err)
			return nil, err
		}
		existingLocale = Locale(allergenLocale)
	case Feature:
		featureLocale := FeaturesLocale{
			Name:     locale.Name,
			Locale:   locale.Locale,
			// ParentID: locale.ParentID,
		}
		err = db.Where(&featureLocale).FirstOrInit(&featureLocale).Error
		if err != nil {
			fmt.Println("Error finding/creating Feature locale:", err)
			return nil, err
		}
		existingLocale = Locale(featureLocale)
	case Nutrient:
		nutrientLocale := NutrientsLocale{
			Name:     locale.Name,
			Locale:   locale.Locale,
			// ParentID: locale.ParentID,
		}
		err = db.Where(&nutrientLocale).FirstOrInit(&nutrientLocale).Error
		if err != nil {
			fmt.Println("Error finding/creating Nutrient locale:", err)
			return nil, err
		}
		existingLocale = Locale(nutrientLocale)
	default:
		return nil, fmt.Errorf("unknown entity type")
	}

	return &existingLocale, nil
}


func insertEntityWithLocales[T EntityInterface](db *gorm.DB, entity T, locales *[]Locale) (*T, error) {
	// fmt.Printf("Initial Locales: %+v\n", locales)

	for i, locale := range *locales {
		// Find or create locale for the given entity type
		foundLocale, err := findOrCreateLocale(db, &locale, entity)
		if err != nil {
			return nil, err
		}

		// Replace the locale in the slice with the one from the database (if it's newly created)
		(*locales)[i] = *foundLocale
	}
	// fmt.Printf("Locales: %+v\n", *locales)

	// Check if the first locale has ID, if not, it's new and not in the database yet
	if (*locales)[0].ID == 0 {
		// Create the entity
		if err := db.Create(&entity).Error; err != nil {
			return nil, err
		}
		// fmt.Printf("Creating Entity: %+v\n", entity)
	} else {
		// Get Typed entity
		entity.SetID((*locales)[0].ParentID)
		if err := db.First(&entity).Error; err != nil {
			return nil, err
		}
		// fmt.Printf("Found Entity: %+v\n", entity)
	}

	// Finally, create or find the entity itself
	if err := db.FirstOrCreate(&entity).Error; err != nil {
		return nil, err
	}

	// fmt.Printf("Entity: %+v\n", entity)

	// If locale wasnt created, create it now
	for _, locale := range *locales {
		if locale.ID == 0 {
			locale.ParentID = entity.GetID()
			switch any(entity).(type) {
			case Additive:
				additiveLocale := AdditivesLocale{
					Name:     locale.Name,
					Locale:   locale.Locale,
					ParentID: locale.ParentID,
				}
				err := db.Where(&additiveLocale).Create(&additiveLocale).Error
				if err != nil {
					fmt.Println("Error finding/creating Additive locale:", err)
					return nil, err
				}
			case Allergen:
				allergenLocale := AllergensLocale{
					Name:     locale.Name,
					Locale:   locale.Locale,
					ParentID: locale.ParentID,
				}
				err := db.Where(&allergenLocale).Create(&allergenLocale).Error
				if err != nil {
					fmt.Println("Error finding/creating Allergen locale:", err)
					return nil, err
				}
			case Feature:
				featureLocale := FeaturesLocale{
					Name:     locale.Name,
					Locale:   locale.Locale,
					ParentID: locale.ParentID,
				}
				err := db.Where(&featureLocale).Create(&featureLocale).Error
				if err != nil {
					fmt.Println("Error finding/creating Feature locale:", err)
					return nil, err
				}
			case Nutrient:
				nutrientLocale := NutrientsLocale{
					Name:     locale.Name,
					Locale:   locale.Locale,
					ParentID: locale.ParentID,
				}
				err := db.Where(&nutrientLocale).Create(&nutrientLocale).Error
				if err != nil {
					fmt.Println("Error finding/creating Nutrient locale:", err)
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unknown entity type")
			}
		}
	}

	return &entity, nil
}



// Inserters
