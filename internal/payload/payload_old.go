package payload

import (
	"fmt"
	"mensa-queue/internal/repository"
	"time"

	"gorm.io/gorm"
)

type Mensa int32

const (
	NeuesPalais      Mensa = 9600
	Griebnitzsee     Mensa = 9601
	Golm             Mensa = 9602
	Filmuniversitaet Mensa = 9603
	FHP              Mensa = 9604
	Wildau           Mensa = 9605
	Brandenburg      Mensa = 9606
)

type Language int

const (
	DE Language = iota + 1
	EN
)

func (l Language) String() string {
	return [...]string{"", "de", "en"}[l]
}

type LocalizedString map[repository.EnumLocaleLocale]string

type LocalizedValue[T any] struct {
	DE T
	EN T
}

type Recipe struct {
	ID             uint `gorm:"primaryKey"`
	PriceStudents  *float64
	PriceEmployees *float64
	PriceGuests    *float64
	MensaProvider  int64 `gorm:"column:mensa_provider_id"`
	AIThumbnailID  uint  `gorm:"column:ai_thumbnail_id"`
}

func (r Recipe) GetID() uint   { return r.ID }
func (r Recipe) SetID(id uint) { r.ID = id }

type RecipesRel struct {
	ID          uint `gorm:"primaryKey"`
	ParentID    uint `gorm:"column:parent_id"`
	Path        string
	AdditivesID *uint `gorm:"column:additives_id"`
	AllergensID *uint `gorm:"column:allergens_id"`
	FeaturesID  *uint `gorm:"column:features_id"`
}

type Locale struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	Locale   string `gorm:"column:_locale"`
	ParentID uint   `gorm:"column:_parent_id"`
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
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	Locale   string `gorm:"column:_locale"`
	ParentID uint   `gorm:"column:_parent_id"`
}

// Implement LocaleInterface methods for each locale type
func (l AdditivesLocale) GetID() uint          { return l.ID }
func (l AdditivesLocale) GetName() string      { return l.Name }
func (l AdditivesLocale) GetLocale() string    { return l.Locale }
func (l AdditivesLocale) GetParentID() uint    { return l.ParentID }
func (l *AdditivesLocale) SetParentID(id uint) { l.ParentID = id }

type RecipesLocale struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	Locale   string `gorm:"column:_locale"`
	ParentID uint   `gorm:"column:_parent_id"`
}

func (l RecipesLocale) GetID() uint          { return l.ID }
func (l RecipesLocale) GetName() string      { return l.Name }
func (l RecipesLocale) GetLocale() string    { return l.Locale }
func (l RecipesLocale) GetParentID() uint    { return l.ParentID }
func (l *RecipesLocale) SetParentID(id uint) { l.ParentID = id }

type AllergensLocale struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	Locale   string `gorm:"column:_locale"`
	ParentID uint   `gorm:"column:_parent_id"`
}

func (l AllergensLocale) GetID() uint          { return l.ID }
func (l AllergensLocale) GetName() string      { return l.Name }
func (l AllergensLocale) GetLocale() string    { return l.Locale }
func (l AllergensLocale) GetParentID() uint    { return l.ParentID }
func (l *AllergensLocale) SetParentID(id uint) { l.ParentID = id }

type FeaturesLocale struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	Locale   string `gorm:"column:_locale"`
	ParentID uint   `gorm:"column:_parent_id"`
}

func (l FeaturesLocale) GetID() uint          { return l.ID }
func (l FeaturesLocale) GetName() string      { return l.Name }
func (l FeaturesLocale) GetLocale() string    { return l.Locale }
func (l FeaturesLocale) GetParentID() uint    { return l.ParentID }
func (l *FeaturesLocale) SetParentID(id uint) { l.ParentID = id }

type NutrientsLocale struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	Locale   string `gorm:"column:_locale"`
	ParentID uint   `gorm:"column:_parent_id"`
}

func (l NutrientsLocale) GetID() uint          { return l.ID }
func (l NutrientsLocale) GetName() string      { return l.Name }
func (l NutrientsLocale) GetLocale() string    { return l.Locale }
func (l NutrientsLocale) GetParentID() uint    { return l.ParentID }
func (l *NutrientsLocale) SetParentID(id uint) { l.ParentID = id }

type Additive struct {
	ID uint `gorm:"primaryKey"`
}

func (a Additive) GetID() uint   { return a.ID }
func (a Additive) SetID(id uint) { a.ID = id }

type Allergen struct {
	ID uint `gorm:"primaryKey"`
}

func (a Allergen) GetID() uint   { return a.ID }
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

func (n NutrientLabel) GetID() uint   { return n.ID }
func (n NutrientLabel) SetID(id uint) { n.ID = id }

type NutrientUnit struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type Serving struct {
	ID       uint `gorm:"primaryKey"`
	Date     time.Time
	MensaID  uint `gorm:"column:mensa_id"`
	RecipeID uint `gorm:"column:recipe_id"`
}

type NutrientLabelsLocale struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	ParentID uint   `gorm:"column:_parent_id"`
	Locale   string `gorm:"column:_locale"`
}

type Feature struct {
	ID              uint `gorm:"primaryKey"`
	MensaProviderID uint `gorm:"column:mensa_provider_id"`
}

func (f Feature) GetID() uint   { return f.ID }
func (f Feature) SetID(id uint) { f.ID = id }

// findOrCreateLocale abstracts the repeated logic of checking and creating locales.
func findOrCreateLocale[T any](db *gorm.DB, locale *Locale, entity T) (*Locale, error) {
	var existingLocale Locale
	var err error

	switch any(entity).(type) {
	case Additive:
		additiveLocale := AdditivesLocale{
			Name:   locale.Name,
			Locale: locale.Locale,
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
			Name:   locale.Name,
			Locale: locale.Locale,
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
			Name:   locale.Name,
			Locale: locale.Locale,
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
			Name:   locale.Name,
			Locale: locale.Locale,
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
