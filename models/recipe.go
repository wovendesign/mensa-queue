package models

import (
	"mensa-queue/internal/repository"
	"time"
)

type Recipe struct {
	ID              uint
	PriceStudents   *float64
	PriceEmployees  *float64
	PriceGuests     *float64
	MensaProviderID *int32
	AIThumbnailID   uint
	Localization    *RecipeLocalization
	Serving         *Serving
}

type RecipeLocalization struct {
	Locales   []*repository.InsertLocaleParams
	Allergen  [][]*repository.InsertLocaleParams
	Additives [][]*repository.InsertLocaleParams
	Features  [][]*repository.InsertLocaleParams
	Nutrients []*NutrientLocalization
	Category  repository.EnumRecipesCategory
}

type NutrientLocalization struct {
	Unit    string
	Value   float64
	Locales []*repository.InsertLocaleParams
}

type Serving struct {
	MensaID *int32
	Date    time.Time
}

type StringLocalization map[repository.EnumLocaleLocale]string
