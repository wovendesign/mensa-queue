package parsers

import (
	"encoding/json"
	"fmt"
	"io"
	"mensa-queue/internal/payload"

	"net/http"
	"strconv"
	"strings"
)

type Model string
const (
	AdditivesModel Model = "additives"
	AllergensModel Model = "allergens"
	FeaturesModel Model = "features"
	FoodModel Model = "menu"
	CategoryModel Model = "mealCategory"
)

func sendRequestToSWT(model Model, mensa payload.Mensa, languageType payload.Language) ([]byte, error) {
	client := &http.Client{}
	url := "https://swp.webspeiseplan.de/index.php?token=55ed21609e26bbf68ba2b19390bf7961"
	reqURL := fmt.Sprintf("%s&model=%s&location=%d&languagetype=%d", url, model, mensa, languageType)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Referer", "https://swp.webspeiseplan.de/InitialConfig")


	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

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

func ParseModel[T any](model Model, mensa payload.Mensa, languageType payload.Language) (*SWTResponse[T], error) {
    body, err := sendRequestToSWT(model, mensa, languageType)
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
	Success bool         `json:"success"`
	Content []T `json:"content"`
}

type FoodContent struct {
	SpeiseplanGerichtData  []SpeiseplanGerichtDatum `json:"speiseplanGerichtData"`
}

type SpeiseplanGerichtDatum struct {
	SpeiseplanAdvancedGericht SpeiseplanAdvancedGericht `json:"speiseplanAdvancedGericht"`
	Zusatzinformationen        Zusatzinformationen         `json:"zusatzinformationen"`
	AllergenIDsString              string                     `json:"allergeneIds"`
	AdditivesIDsString           *string                    `json:"zusatzstoffeIds,omitempty"`
	FeaturesIDsString        string                     `json:"gerichtmerkmaleIds"`
}

type SpeiseplanAdvancedGericht struct {
	ID                          int64  `json:"id"`
	Active                       bool   `json:"aktiv"`
	Date                       string `json:"datum"`
	RecipeCategoryID          int64  `json:"gerichtkategorieID"`
	RecipeName                 string `json:"gerichtname"`
	ZusatzinformationenID        int64  `json:"zusatzinformationenID"`
	SpeiseplanAdvancedID        int64  `json:"speiseplanAdvancedID"`
	TimestampLog                string `json:"timestampLog"`
	UserID                 int64  `json:"benutzerID"`
}

type Zusatzinformationen struct {
	ID                         int64           `json:"id"`
	GerichtnameAlternative     string          `json:"gerichtnameAlternative"`
	MitarbeiterpreisDecimal2   float64         `json:"mitarbeiterpreisDecimal2"`
	GaestepreisDecimal2       float64         `json:"gaestepreisDecimal2"`
	EnaehrungsampelID         *json.RawMessage `json:"ernaehrungsampelID,omitempty"`
	NwkjInteger                int64           `json:"nwkjInteger"`
	NwkcalInteger              int64           `json:"nwkcalInteger"`
	NwfettDecimal1             float64         `json:"nwfettDecimal1"`
	NwfettsaeurenDecimal1      float64         `json:"nwfettsaeurenDecimal1"`
	NwkohlehydrateDecimal1     float64         `json:"nwKohlehydrateDecimal1"`
	NwzuckerDecimal1           float64         `json:"nwzuckerDecimal1"`
	NweiweissDecimal1          float64         `json:"nweiweissDecimal1"`
	NwsalzDecimal1             float64         `json:"nwsalzDecimal1"`
	NwbeDecimal2               *json.RawMessage `json:"nwbeDecimal2,omitempty"`
	AllowFeedback              bool            `json:"allowFeedback"`
	GerichtImage               *json.RawMessage `json:"gerichtImage,omitempty"`
	Lieferanteninfo            *json.RawMessage `json:"lieferanteninfo,omitempty"`
	LieferanteninfoLink        *json.RawMessage `json:"lieferanteninfoLink,omitempty"`
	EdFaktorDecimal1           *json.RawMessage `json:"edFaktorDecimal1,omitempty"`
	Plu                        *string          `json:"plu,omitempty"`
	Price3Decimal2            float64         `json:"price3Decimal2"`
	Price4Decimal2            *json.RawMessage `json:"price4Decimal2,omitempty"`
	Contingent                *json.RawMessage `json:"contingent,omitempty"`
	TaxRateDecimal2           *json.RawMessage `json:"taxRateDecimal2,omitempty"`
	IngredientList            *json.RawMessage `json:"ingredientList,omitempty"`
	Sustainability             Sustainability    `json:"sustainability"`
}

type Sustainability struct {
	CO2        *Co2           `json:"co2,omitempty"`
	Nutriscore *json.RawMessage `json:"nutriscore,omitempty"`
	TrafficLight *json.RawMessage `json:"trafficLight,omitempty"`
}

type Co2 struct {
	ID        int64 `json:"id"`
	CO2Value  int64 `json:"co2Value"`
	Unit      Unit  `json:"unit"`
}

type Unit string

const (
	UnitG Unit = "g"
)

type AdditivesResponse struct {
	Success bool               `json:"success"`
	Content []AdditivesContent `json:"content"`
}

type AdditivesContent struct {
	ID              int64         `json:"id"`
	Name            string        `json:"name"`
	Kuerzel         string        `json:"kuerzel"`
	Beschreibung    json.RawMessage `json:"beschreibung"`
	ZusatzstoffeID  int64         `json:"zusatzstoffeID"`
	LanguageTypeID  int64         `json:"languageTypeID"`
	BenutzerID      int64         `json:"benutzerID"`
}

type AllergensResponse struct {
	Success bool               `json:"success"`
	Content []AllergensContent `json:"content"`
}

type AllergensContent struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	AllergeneID  int64         `json:"allergeneID"`
}


type FeatureContent struct {
	ID                     int64         `json:"id"`
	Name                   string        `json:"name"`
	NameAlternative        json.RawMessage `json:"nameAlternative"`
	Kuerzel                string        `json:"kuerzel"`
	LogoImage              *string       `json:"logoImage,omitempty"`
	RGBColor               json.RawMessage `json:"rgbColor"`
	ReihenfolgeInApp       int64         `json:"reihenfolgeInApp"`
	ShowInSpeiseplanOverview bool          `json:"showInSpeiseplanOverview"`
	ShowNotInFilter        bool          `json:"showNotInFilter"`
	Beschreibung           json.RawMessage `json:"beschreibung"`
	GerichtmerkmalID       int64         `json:"gerichtmerkmalID"`
	LanguageTypeID         int64         `json:"languageTypeID"`
	TimestampLog           string        `json:"timestampLog"`
	BenutzerID             int64         `json:"benutzerID"`
}

type MensaRequest struct {
	Food      []FoodContent      `json:"food"`
	Additives []AdditivesContent  `json:"additives"`
	Allergens []AllergensContent  `json:"allergens"`
	Features  []FeatureContent    `json:"features"`
}

type Recipe struct {
	Title    string   `json:"title"`
	Diet     Diet     `json:"diet"`
	Prices   Prices   `json:"prices"`
	Nutrients Nutrients `json:"nutrients"`
}

type Diet string

const (
	DietVegan       Diet = "Vegan"
	DietVegetarian  Diet = "Vegetarian"
	DietMeat        Diet = "Meat"
	DietFish        Diet = "Fish"
)

type Prices struct {
	Student  *float64 `json:"student,omitempty"`
	Employee *float64 `json:"employee,omitempty"`
	Other    *float64 `json:"other,omitempty"`
}

type Nutrients struct {
	Calories          float64 `json:"calories"`
	Protein           float64 `json:"protein"`
	Fat               float64 `json:"fat"`
	SaturatedFat      float64 `json:"saturatedFat"`
	Carbs             float64 `json:"carbs"`
	Sugar             float64 `json:"sugar"`
	Salt              float64 `json:"salt"`
}

type Features string

const (
	FeatureVegetarian     Features = "Vegetarian"
	FeatureVegan          Features = "Vegan"
	FeaturePoultry        Features = "Poultry"
	FeatureBeef           Features = "Beef"
	FeaturePork           Features = "Pork"
	FeatureGame           Features = "Game"
	FeatureRegional       Features = "Regional"
	FeatureAlcohol        Features = "Alcohol"
	FeatureFish           Features = "Fish"
	FeatureLamb           Features = "Lamb"
	FeatureGarlic         Features = "Garlic"
	FeatureBear           Features = "Bear"
	FeatureLeek           Features = "Leek"
	FeatureRennet         Features = "Rennet"
)

type FeatureList struct {
	Vegetarian       []int64 `json:"vegetarian"`
	Vegan            []int64 `json:"vegan"`
	Poultry          []int64 `json:"poultry"`
		Beef             []int64 `json:"beef"`
		Pork             []int64 `json:"pork"`
		Game             []int64 `json:"game"`
		Regional         []int64 `json:"regional"`
		Alcohol          []int64 `json:"alcohol"`
		Fish             []int64 `json:"fish"`
		Lamb             []int64 `json:"lamb"`
		GarlicBearLeek   []int64 `json:"garlicBearLeek"`
		Rennet           []int64 `json:"rennet"`
}


func ParsePotsdamMensaData() (*[]FoodContent, error) {
    body, err := sendRequestToSWT(FoodModel, payload.NeuesPalais, payload.DE)
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

	return &foodResponse.Content, nil
}

type MealCategoryResponse struct {
	Success bool `json:"success"`
	Content []MealCategory `json:"content"`
}

type MealCategory struct {
	ID int64 `json:"gerichtkategorieID"`
	Name string `json:"name"`
	LanguageType int64 `json:"languageTypeID"`
}

func ParseMealCategory() (*[]MealCategory, error) {
	body, err := sendRequestToSWT(CategoryModel, payload.NeuesPalais, payload.DE)

	var mealCategoryResponse MealCategoryResponse

	err = json.Unmarshal(body, &mealCategoryResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Check if the API indicates failure
	if !mealCategoryResponse.Success {
		return nil, fmt.Errorf("API error: success is false, content: %v", mealCategoryResponse.Content)
	}

	return &mealCategoryResponse.Content, nil
}

func ExtractNutrients(food SpeiseplanGerichtDatum) (*[]payload.LocalNutrient, error) {
	nutrients := make([]payload.LocalNutrient, 0)

	nutrients = append(nutrients, payload.LocalNutrient{
		Nutrient: payload.Nutrient{
			NutrientLabel: payload.NutrientLabel{
				Unit: payload.NutrientUnit{
					Name: "kcal",
				},
			},
			NutrientValue: payload.NutrientValue{
				Value: float64(food.Zusatzinformationen.NwkcalInteger),
			},
		},
		Name: payload.LocalizedString {
			payload.DE: "Kalorien",
			payload.EN: "Calories",
		},
	})

	nutrients = append(nutrients, payload.LocalNutrient{
		Nutrient: payload.Nutrient{
			NutrientLabel: payload.NutrientLabel{
				Unit: payload.NutrientUnit{
					Name: "g",
				},
			},
			NutrientValue: payload.NutrientValue{
				Value: float64(food.Zusatzinformationen.NwfettDecimal1),
			},
		},
		Name: payload.LocalizedString {
			payload.DE: "Fett",
			payload.EN: "Fat",
		},
	})

	nutrients = append(nutrients, payload.LocalNutrient{
		Nutrient: payload.Nutrient{
			NutrientLabel: payload.NutrientLabel{
				Unit: payload.NutrientUnit{
					Name: "g",
				},
			},
			NutrientValue: payload.NutrientValue{
				Value: float64(food.Zusatzinformationen.NwfettsaeurenDecimal1),
			},
		},
		Name: payload.LocalizedString {
			payload.DE: "Gesättigte Fettsäuren",
			payload.EN: "Saturated Fatty Acids",
		},
	})

	nutrients = append(nutrients, payload.LocalNutrient{
		Nutrient: payload.Nutrient{
			NutrientLabel: payload.NutrientLabel{
				Unit: payload.NutrientUnit{
					Name: "g",
				},
			},
			NutrientValue: payload.NutrientValue{
				Value: float64(food.Zusatzinformationen.NwkohlehydrateDecimal1),
			},
		},
		Name: payload.LocalizedString {
			payload.DE: "Kohlenhydrate",
			payload.EN: "Carbohydrates",
		},
	})

	nutrients = append(nutrients, payload.LocalNutrient{
		Nutrient: payload.Nutrient{
			NutrientLabel: payload.NutrientLabel{
				Unit: payload.NutrientUnit{
					Name: "g",
				},
			},
			NutrientValue: payload.NutrientValue{
				Value: float64(food.Zusatzinformationen.NwzuckerDecimal1),
			},
		},
		Name: payload.LocalizedString {
			payload.DE: "Zucker",
			payload.EN: "Sugar",
		},
	})

	nutrients = append(nutrients, payload.LocalNutrient{
		Nutrient: payload.Nutrient{
			NutrientLabel: payload.NutrientLabel{
				Unit: payload.NutrientUnit{
					Name: "g",
				},
			},
			NutrientValue: payload.NutrientValue{
				Value: float64(food.Zusatzinformationen.NweiweissDecimal1),
			},
		},
		Name: payload.LocalizedString {
			payload.DE: "Eiweiß",
			payload.EN: "Protein",
		},
	})

	nutrients = append(nutrients, payload.LocalNutrient{
		Nutrient: payload.Nutrient{
			NutrientLabel: payload.NutrientLabel{
				Unit: payload.NutrientUnit{
					Name: "g",
				},
			},
			NutrientValue: payload.NutrientValue{
				Value: float64(food.Zusatzinformationen.NwsalzDecimal1),
			},
		},
		Name: payload.LocalizedString {
			payload.DE: "Salz",
			payload.EN: "Salt",
		},
	})

	return &nutrients, nil
}


type AdditiveResponse struct {
	ID int64 `json:"zusatzstoffeID"`
	Name string `json:"name"`
}

func ParseAdditives(languages []payload.Language) (map[int64]payload.LocalizedString, error) {
	allAdditives := make(map[int64]payload.LocalizedString)
	for _, language := range languages {
		additives, err := ParseModel[AdditiveResponse](AdditivesModel, payload.NeuesPalais, language)
		if err != nil {
			return nil, err
		}
		for _, additive := range additives.Content {
			// Check if the additive ID already exists in the map
			localizedString, exists := allAdditives[additive.ID]
			if !exists {
				// If it doesn't exist, initialize the map
				localizedString = payload.LocalizedString{}
			}

			// Add or update the name for the current language
			localizedString[language] = additive.Name

			// Update the map with the modified localized string
			allAdditives[additive.ID] = localizedString
		}
	}

	return allAdditives, nil
}

func ExtractAdditives(food SpeiseplanGerichtDatum, additives map[int64]payload.LocalizedString, languages []payload.Language) (*[][]payload.Locale, error) {
	if food.AdditivesIDsString == nil || len(*food.AdditivesIDsString) == 0 {
		return nil, nil
	}
	additivesArray := strings.Split(*food.AdditivesIDsString, ",")

	var result [][]payload.Locale

	for _, additiveID := range additivesArray {
		additiveIDInt, err := strconv.Atoi(additiveID)
		if err != nil {
			return nil, fmt.Errorf("error parsing additive ID: %w", err)
		}

		var additivePair []payload.Locale
		for _, language := range languages {
			additivePair = append(additivePair, payload.Locale{
				Name:   additives[int64(additiveIDInt)][language],
				Locale: language.String(),
			})
		}

		result = append(result, additivePair)
	}

	return &result, nil
}

type AllergenResponse struct {
	ID int64 `json:"allergeneID"`
	Name string `json:"name"`
}

func ParseAllergens(languages []payload.Language) (map[int64]payload.LocalizedString, error) {
	allAllergens := make(map[int64]payload.LocalizedString)

	for _, language := range languages {
		allergens, err := ParseModel[AllergenResponse](AllergensModel, payload.NeuesPalais, language)
		if err != nil {
			return nil, err
		}

		for _, allergen := range allergens.Content {
			// Check if the allergen ID already exists in the map
			localizedString, exists := allAllergens[allergen.ID]
			if !exists {
				// If it doesn't exist, initialize the map
				localizedString = payload.LocalizedString{}
			}

			// Add or update the name for the current language
			localizedString[language] = allergen.Name

			// Update the map with the modified localized string
			allAllergens[allergen.ID] = localizedString
		}
	}

	return allAllergens, nil
}

func ExtractAllergens(food SpeiseplanGerichtDatum, allergens map[int64]payload.LocalizedString, languages []payload.Language) (*[][]payload.Locale, error) {
	if food.AllergenIDsString == "" {
		return nil, nil
	}

	allergenIDs := strings.Split(food.AllergenIDsString, ",")
	var result [][]payload.Locale

	for _, allergenID := range allergenIDs {
		allergenIDInt, err := strconv.Atoi(allergenID)
		if err != nil {
			return nil, fmt.Errorf("error parsing allergen ID: %w", err)
		}

		var allergenPair []payload.Locale
		for _, language := range languages {
			allergenPair = append(allergenPair, payload.Locale{
				Name:   allergens[int64(allergenIDInt)][language],
				Locale: language.String(),
			})
		}

		result = append(result, allergenPair)
	}

	return &result, nil
}

type FeatureResponse struct {
	ID int64 `json:"gerichtmerkmalID"`
	Name string `json:"name"`
}

func ParseFeatures(languages []payload.Language) (map[int64]payload.LocalizedString, error) {
	allFeatures := make(map[int64]payload.LocalizedString)
	for _, language := range languages {
		features, err := ParseModel[FeatureResponse](FeaturesModel, payload.NeuesPalais, language)
		if err != nil {
			return nil, err
		}
		for _, feature := range features.Content {
			// Check if the feature ID already exists in the map
			localizedString, exists := allFeatures[feature.ID]
			if !exists {
				// If it doesn't exist, initialize the map
				localizedString = payload.LocalizedString{}
			}

			// Add or update the name for the current language
			localizedString[language] = feature.Name

			// Update the map with the modified localized string
			allFeatures[feature.ID] = localizedString
		}
	}

	return allFeatures, nil
}

func ExtractFeatures(food SpeiseplanGerichtDatum, features map[int64]payload.LocalizedString, languages []payload.Language) (*[][]payload.Locale, error) {
	if food.FeaturesIDsString == "" {
		return nil, nil
	}

	featureIDs := strings.Split(food.FeaturesIDsString, ",")
	var result [][]payload.Locale

	for _, featureID := range featureIDs {
		featureIDInt, err := strconv.Atoi(featureID)
		if err != nil {
			return nil, fmt.Errorf("error parsing allergen ID: %w", err)
		}

		var featurePair []payload.Locale
		for _, language := range languages {
			featurePair = append(featurePair, payload.Locale{
				Name:   features[int64(featureIDInt)][language],
				Locale: language.String(),
			})
		}

		result = append(result, featurePair)
	}

	return &result, nil
}
