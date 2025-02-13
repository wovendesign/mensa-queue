package stw_brandenburg_west

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mensa-queue/internal/repository"
	"mensa-queue/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *StwBrandenburgWestMensa) ParseMenu() ([]*models.Recipe, error) {
	languages := []repository.EnumLocaleLocale{repository.EnumLocaleLocaleDe, repository.EnumLocaleLocaleEn}
	foodContent, err := ParsePotsdamMensaData(&m.StwID)
	if err != nil {
		log.Fatal(err)
	}

	additiveMap, err := ParseAdditives(languages, &m.StwID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	allergenMap, err := ParseAllergens(languages, &m.StwID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	featureMap, err := ParseFeatures(languages, &m.StwID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	categoryMap, err := ParseMealCategory(&m.StwID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var recipes []*models.Recipe

	for _, week := range foodContent {
		for _, food := range week.SpeiseplanGerichtData {
			if food.Zusatzinformationen.MitarbeiterpreisDecimal2 == 0 || strings.Contains(food.SpeiseplanAdvancedGericht.RecipeName, "Preis pro") {
				continue
			}

			category, err := ExtractCategories(food, categoryMap)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			nutrients, err := ExtractNutrients(food)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			additives, err := ExtractAdditives(food, additiveMap, languages)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			allergens, err := ExtractAllergens(food, allergenMap, languages)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			features, err := ExtractFeatures(food, featureMap, languages)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			t, err := time.Parse(time.RFC3339, food.SpeiseplanAdvancedGericht.Date)
			t = t.UTC()
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return nil, err
			}

			recipe := models.Recipe{
				PriceStudents:   &food.Zusatzinformationen.MitarbeiterpreisDecimal2,
				PriceGuests:     &food.Zusatzinformationen.GaestepreisDecimal2,
				MensaProviderID: m.Provider.MensaHubID,
				Localization: &models.RecipeLocalization{
					Locales: []*repository.InsertLocaleParams{
						{
							Name:   food.SpeiseplanAdvancedGericht.RecipeName,
							Locale: repository.EnumLocaleLocaleDe,
						},
						{
							Name:   food.Zusatzinformationen.GerichtnameAlternative,
							Locale: repository.EnumLocaleLocaleEn,
						},
					},

					Nutrients: nutrients,
					Allergen:  allergens,
					Additives: additives,
					Features:  features,
					Category:  category,
				},
				Serving: &models.Serving{
					MensaID: m.MensaHubID,
					Date:    t,
				},
			}
			recipes = append(recipes, &recipe)
		}
	}
	return recipes, nil
}

type Model string

const (
	AdditivesModel Model = "additives"
	AllergensModel Model = "allergens"
	FeaturesModel  Model = "features"
	FoodModel      Model = "menu"
	CategoryModel  Model = "mealCategory"
)

func sendRequestToSWT(model Model, stwMensaId *int32, languageType repository.EnumLocaleLocale) ([]byte, error) {
	client := &http.Client{}
	url := "https://swp.webspeiseplan.de/index.php?token=55ed21609e26bbf68ba2b19390bf7961"

	var languageInt int
	switch languageType {
	case repository.EnumLocaleLocaleDe:
		languageInt = 1
	case repository.EnumLocaleLocaleEn:
		languageInt = 2
	}

	reqURL := fmt.Sprintf("%s&model=%s&location=%d&languagetype=%d", url, model, *stwMensaId, languageInt)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Referer", "https://swp.webspeiseplan.de/InitialConfig")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
			return
		}
	}(resp.Body)

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func ParseModel[T any](model Model, stwMensaId *int32, languageType repository.EnumLocaleLocale) (*SWTResponse[T], error) {
	body, err := sendRequestToSWT(model, stwMensaId, languageType)
	if err != nil {
		return nil, err
	}

	var response SWTResponse[T]
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &response, nil
}

type SWTResponse[T any] struct {
	Success bool `json:"success"`
	Content []T  `json:"content"`
}

type FoodContent struct {
	SpeiseplanGerichtData []SpeiseplanGerichtDatum `json:"speiseplanGerichtData"`
}

type SpeiseplanGerichtDatum struct {
	SpeiseplanAdvancedGericht SpeiseplanAdvancedGericht `json:"speiseplanAdvancedGericht"`
	Zusatzinformationen       Zusatzinformationen       `json:"zusatzinformationen"`
	AllergenIDsString         string                    `json:"allergeneIds"`
	AdditivesIDsString        *string                   `json:"zusatzstoffeIds,omitempty"`
	FeaturesIDsString         string                    `json:"gerichtmerkmaleIds"`
}

type SpeiseplanAdvancedGericht struct {
	ID                    int64  `json:"id"`
	Active                bool   `json:"aktiv"`
	Date                  string `json:"datum"`
	RecipeCategoryID      int64  `json:"gerichtkategorieID"`
	RecipeName            string `json:"gerichtname"`
	ZusatzinformationenID int64  `json:"zusatzinformationenID"`
	SpeiseplanAdvancedID  int64  `json:"speiseplanAdvancedID"`
	TimestampLog          string `json:"timestampLog"`
	UserID                int64  `json:"benutzerID"`
}

type Zusatzinformationen struct {
	ID                       int64            `json:"id"`
	GerichtnameAlternative   string           `json:"gerichtnameAlternative"`
	MitarbeiterpreisDecimal2 float64          `json:"mitarbeiterpreisDecimal2"`
	GaestepreisDecimal2      float64          `json:"gaestepreisDecimal2"`
	EnaehrungsampelID        *json.RawMessage `json:"ernaehrungsampelID,omitempty"`
	NwkjInteger              int64            `json:"nwkjInteger"`
	NwkcalInteger            int64            `json:"nwkcalInteger"`
	NwfettDecimal1           float64          `json:"nwfettDecimal1"`
	NwfettsaeurenDecimal1    float64          `json:"nwfettsaeurenDecimal1"`
	NwkohlehydrateDecimal1   float64          `json:"nwKohlehydrateDecimal1"`
	NwzuckerDecimal1         float64          `json:"nwzuckerDecimal1"`
	NweiweissDecimal1        float64          `json:"nweiweissDecimal1"`
	NwsalzDecimal1           float64          `json:"nwsalzDecimal1"`
	NwbeDecimal2             *json.RawMessage `json:"nwbeDecimal2,omitempty"`
	AllowFeedback            bool             `json:"allowFeedback"`
	GerichtImage             *json.RawMessage `json:"gerichtImage,omitempty"`
	Lieferanteninfo          *json.RawMessage `json:"lieferanteninfo,omitempty"`
	LieferanteninfoLink      *json.RawMessage `json:"lieferanteninfoLink,omitempty"`
	EdFaktorDecimal1         *json.RawMessage `json:"edFaktorDecimal1,omitempty"`
	Plu                      *string          `json:"plu,omitempty"`
	Price3Decimal2           float64          `json:"price3Decimal2"`
	Price4Decimal2           *json.RawMessage `json:"price4Decimal2,omitempty"`
	Contingent               *json.RawMessage `json:"contingent,omitempty"`
	TaxRateDecimal2          *json.RawMessage `json:"taxRateDecimal2,omitempty"`
	IngredientList           *json.RawMessage `json:"ingredientList,omitempty"`
	Sustainability           Sustainability   `json:"sustainability"`
}

type Sustainability struct {
	CO2          *Co2             `json:"co2,omitempty"`
	Nutriscore   *json.RawMessage `json:"nutriscore,omitempty"`
	TrafficLight *json.RawMessage `json:"trafficLight,omitempty"`
}

type Co2 struct {
	ID       int64 `json:"id"`
	CO2Value int64 `json:"co2Value"`
	Unit     Unit  `json:"unit"`
}

type Unit string

func ParsePotsdamMensaData(stwMensaId *int32) ([]FoodContent, error) {
	body, err := sendRequestToSWT(FoodModel, stwMensaId, repository.EnumLocaleLocaleDe)
	if err != nil {
		return nil, err
	}

	var foodResponse SWTResponse[FoodContent]
	// Parse the JSON data into the struct
	err = json.Unmarshal(body, &foodResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return foodResponse.Content, nil
}

type MealCategoryResponse struct {
	Success bool           `json:"success"`
	Content []MealCategory `json:"content"`
}

type MealCategory struct {
	ID           int64  `json:"gerichtkategorieID"`
	Name         string `json:"name"`
	LanguageType int64  `json:"languageTypeID"`
}

func ParseMealCategory(stwMensaId *int32) (map[int64]MealCategory, error) {
	body, err := sendRequestToSWT(CategoryModel, stwMensaId, repository.EnumLocaleLocaleDe)

	var mealCategoryResponse MealCategoryResponse

	err = json.Unmarshal(body, &mealCategoryResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Check if the API indicates failure
	if !mealCategoryResponse.Success {
		return nil, fmt.Errorf("API error: success is false, content: %v", mealCategoryResponse.Content)
	}
	content := mealCategoryResponse.Content

	categories := make(map[int64]MealCategory, len(content))

	for _, category := range content {
		categories[category.ID] = category
	}

	return categories, nil
}

func ExtractCategories(food SpeiseplanGerichtDatum, categoryMap map[int64]MealCategory) (repository.EnumRecipesCategory, error) {
	category := categoryMap[food.SpeiseplanAdvancedGericht.RecipeCategoryID]

	if strings.Contains(category.Name, "Dessert") {
		return repository.EnumRecipesCategoryDessert, nil
	}

	return repository.EnumRecipesCategoryMain, nil
}

func ExtractNutrients(food SpeiseplanGerichtDatum) ([]*models.NutrientLocalization, error) {
	nutrients := make([]*models.NutrientLocalization, 0)

	nutrients = append(nutrients, &models.NutrientLocalization{
		Unit:  "kcal",
		Value: float64(food.Zusatzinformationen.NwkcalInteger),
		Locales: []*repository.InsertLocaleParams{
			{
				Name:   "Kalorien",
				Locale: repository.EnumLocaleLocaleDe,
			},
			{
				Name:   "Calories",
				Locale: repository.EnumLocaleLocaleEn,
			},
		},
	})

	nutrients = append(nutrients, &models.NutrientLocalization{
		Unit:  "g",
		Value: food.Zusatzinformationen.NwfettDecimal1,
		Locales: []*repository.InsertLocaleParams{
			{
				Name:   "Fett",
				Locale: repository.EnumLocaleLocaleDe,
			},
			{
				Name:   "Fat",
				Locale: repository.EnumLocaleLocaleEn,
			},
		},
	})

	nutrients = append(nutrients, &models.NutrientLocalization{
		Unit:  "g",
		Value: food.Zusatzinformationen.NwfettsaeurenDecimal1,
		Locales: []*repository.InsertLocaleParams{
			{
				Name:   "Gesättigte Fettsäuren",
				Locale: repository.EnumLocaleLocaleDe,
			},
			{
				Name:   "Saturated Fatty Acids",
				Locale: repository.EnumLocaleLocaleEn,
			},
		},
	})

	nutrients = append(nutrients, &models.NutrientLocalization{
		Unit:  "g",
		Value: food.Zusatzinformationen.NwkohlehydrateDecimal1,
		Locales: []*repository.InsertLocaleParams{
			{
				Name:   "Kohlenhydrate",
				Locale: repository.EnumLocaleLocaleDe,
			},
			{
				Name:   "Carbohydrates",
				Locale: repository.EnumLocaleLocaleEn,
			},
		},
	})

	nutrients = append(nutrients, &models.NutrientLocalization{
		Unit:  "g",
		Value: food.Zusatzinformationen.NwzuckerDecimal1,
		Locales: []*repository.InsertLocaleParams{
			{
				Name:   "Zucker",
				Locale: repository.EnumLocaleLocaleDe,
			},
			{
				Name:   "Sugar",
				Locale: repository.EnumLocaleLocaleEn,
			},
		},
	})

	nutrients = append(nutrients, &models.NutrientLocalization{
		Unit:  "g",
		Value: food.Zusatzinformationen.NweiweissDecimal1,
		Locales: []*repository.InsertLocaleParams{
			{
				Name:   "Eiweiß",
				Locale: repository.EnumLocaleLocaleDe,
			},
			{
				Name:   "Protein",
				Locale: repository.EnumLocaleLocaleEn,
			},
		},
	})

	nutrients = append(nutrients, &models.NutrientLocalization{
		Unit:  "g",
		Value: food.Zusatzinformationen.NwsalzDecimal1,
		Locales: []*repository.InsertLocaleParams{
			{
				Name:   "Salz",
				Locale: repository.EnumLocaleLocaleDe,
			},
			{
				Name:   "Salt",
				Locale: repository.EnumLocaleLocaleEn,
			},
		},
	})

	return nutrients, nil
}

type AdditiveResponse struct {
	ID   int64  `json:"zusatzstoffeID"`
	Name string `json:"name"`
}

func ParseAdditives(languages []repository.EnumLocaleLocale, stwMensaId *int32) (map[int64]models.StringLocalization, error) {
	allAdditives := make(map[int64]models.StringLocalization)
	for _, language := range languages {
		additives, err := ParseModel[AdditiveResponse](AdditivesModel, stwMensaId, language)
		if err != nil {
			return nil, err
		}
		for _, additive := range additives.Content {
			// Check if the additive ID already exists in the map
			localizedString, exists := allAdditives[additive.ID]
			if !exists {
				// If it doesn't exist, initialize the map
				localizedString = models.StringLocalization{}
			}

			// Add or update the name for the current language
			localizedString[language] = additive.Name

			// Update the map with the modified localized string
			allAdditives[additive.ID] = localizedString
		}
	}

	return allAdditives, nil
}

func ExtractAdditives(food SpeiseplanGerichtDatum, additives map[int64]models.StringLocalization, languages []repository.EnumLocaleLocale) ([][]*repository.InsertLocaleParams, error) {
	if food.AdditivesIDsString == nil || len(*food.AdditivesIDsString) == 0 {
		return nil, nil
	}
	additivesArray := strings.Split(*food.AdditivesIDsString, ",")

	var result [][]*repository.InsertLocaleParams

	for _, additiveID := range additivesArray {
		additiveIDInt, err := strconv.Atoi(additiveID)
		if err != nil {
			return nil, fmt.Errorf("error parsing additive ID: %w", err)
		}

		var additivePair []*repository.InsertLocaleParams
		for _, language := range languages {
			additivePair = append(additivePair, &repository.InsertLocaleParams{
				Name:   additives[int64(additiveIDInt)][language],
				Locale: language,
			})
		}

		result = append(result, additivePair)
	}

	return result, nil
}

type AllergenResponse struct {
	ID   int64  `json:"allergeneID"`
	Name string `json:"name"`
}

func ParseAllergens(languages []repository.EnumLocaleLocale, stwMensaId *int32) (map[int64]models.StringLocalization, error) {
	allAllergens := make(map[int64]models.StringLocalization)

	for _, language := range languages {
		allergens, err := ParseModel[AllergenResponse](AllergensModel, stwMensaId, language)
		if err != nil {
			return nil, err
		}

		for _, allergen := range allergens.Content {
			// Check if the allergen ID already exists in the map
			localizedString, exists := allAllergens[allergen.ID]
			if !exists {
				// If it doesn't exist, initialize the map
				localizedString = models.StringLocalization{}
			}

			// Add or update the name for the current language
			localizedString[language] = allergen.Name

			// Update the map with the modified localized string
			allAllergens[allergen.ID] = localizedString
		}
	}

	return allAllergens, nil
}

func ExtractAllergens(food SpeiseplanGerichtDatum, allergens map[int64]models.StringLocalization, languages []repository.EnumLocaleLocale) ([][]*repository.InsertLocaleParams, error) {
	if food.AllergenIDsString == "" {
		return nil, nil
	}

	allergenIDs := strings.Split(food.AllergenIDsString, ",")
	var result [][]*repository.InsertLocaleParams

	for _, allergenID := range allergenIDs {
		allergenIDInt, err := strconv.Atoi(allergenID)
		if err != nil {
			return nil, fmt.Errorf("error parsing allergen ID: %w", err)
		}

		var allergenPair []*repository.InsertLocaleParams
		for _, language := range languages {
			allergenPair = append(allergenPair, &repository.InsertLocaleParams{
				Name:   allergens[int64(allergenIDInt)][language],
				Locale: language,
			})
		}

		result = append(result, allergenPair)
	}

	return result, nil
}

type FeatureResponse struct {
	ID   int64  `json:"gerichtmerkmalID"`
	Name string `json:"name"`
}

func ParseFeatures(languages []repository.EnumLocaleLocale, stwMensaId *int32) (map[int64]models.StringLocalization, error) {
	allFeatures := make(map[int64]models.StringLocalization)
	for _, language := range languages {
		features, err := ParseModel[FeatureResponse](FeaturesModel, stwMensaId, language)
		if err != nil {
			return nil, err
		}
		for _, feature := range features.Content {
			// Check if the feature ID already exists in the map
			localizedString, exists := allFeatures[feature.ID]
			if !exists {
				// If it doesn't exist, initialize the map
				localizedString = models.StringLocalization{}
			}

			// Add or update the name for the current language
			localizedString[language] = feature.Name

			// Update the map with the modified localized string
			allFeatures[feature.ID] = localizedString
		}
	}

	return allFeatures, nil
}

func ExtractFeatures(food SpeiseplanGerichtDatum, features map[int64]models.StringLocalization, languages []repository.EnumLocaleLocale) ([][]*repository.InsertLocaleParams, error) {
	if food.FeaturesIDsString == "" {
		return nil, nil
	}

	featureIDs := strings.Split(food.FeaturesIDsString, ",")
	var result [][]*repository.InsertLocaleParams

	for _, featureID := range featureIDs {
		featureIDInt, err := strconv.Atoi(featureID)
		if err != nil {
			return nil, fmt.Errorf("error parsing allergen ID: %w", err)
		}

		var featurePair []*repository.InsertLocaleParams
		for _, language := range languages {
			featurePair = append(featurePair, &repository.InsertLocaleParams{
				Name:   features[int64(featureIDInt)][language],
				Locale: language,
			})
		}

		result = append(result, featurePair)
	}

	return result, nil
}
