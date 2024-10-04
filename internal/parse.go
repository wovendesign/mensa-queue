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

type Language int
const (
	DE Language = iota + 1
	EN
)

type Mensa int
const (
	NeuesPalais Mensa = 9600
	Golm Mensa = 9601
	Teltow Mensa = 9602
	Griebnitzsee Mensa = 9603
)

type Model string
const (
	AdditivesModel Model = "additives"
	AllergensModel Model = "allergens"
	FeaturesModel Model = "features"
	FoodModel Model = "menu"
	CategoryModel Model = "mealCategory"
)

func sendRequestToSWT(model Model, mensa Mensa, languageType Language) ([]byte, error) {
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

type FoodResponse struct {
	Success bool         `json:"success"`
	Content []FoodContent `json:"content"`
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
	Kuerzel      string        `json:"kuerzel"`
	LogoImage    json.RawMessage `json:"logoImage"`
	AllergeneID  int64         `json:"allergeneID"`
	TimestampLog string        `json:"timestampLog"`
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
    body, err := sendRequestToSWT(FoodModel, NeuesPalais, DE)
    if err != nil {
    	return nil, err
    }

	var foodResponse FoodResponse
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

func GetMealCategory() (*[]MealCategory, error) {
	body, err := sendRequestToSWT(CategoryModel, NeuesPalais, DE)

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

func ExtractNutrients(food SpeiseplanGerichtDatum) (*[]payload.Nutrient, error) {
	nutrients := make([]payload.Nutrient, 0)

	// NwkcalInteger: 923,
 //        NwfettDecimal1: 53.33,
 //        NwfettsaeurenDecimal1: 14.67,
 //        NwkohlehydrateDecimal1: 21.82,
 //        NwzuckerDecimal1: 6.49,
 //        NweiweissDecimal1: 49.58,
 //        NwsalzDecimal1: 7.63,

	nutrients = append(nutrients, payload.Nutrient{
		NutrientLabel: payload.NutrientLabel{
			Name: "Kilokalorien",
			Unit: payload.NutrientUnit{
				Name: "kcal",
			},
		},
		NutrientValue: payload.NutrientValue{
			Value: float64(food.Zusatzinformationen.NwkcalInteger),
		},
	})

	nutrients = append(nutrients, payload.Nutrient{
		NutrientLabel: payload.NutrientLabel{
			Name: "Fett",
			Unit: payload.NutrientUnit{
				Name: "g",
			},
		},
		NutrientValue: payload.NutrientValue{
			Value: float64(food.Zusatzinformationen.NwfettDecimal1),
		},
	})

	nutrients = append(nutrients, payload.Nutrient{
		NutrientLabel: payload.NutrientLabel{
			Name: "Fettsäuren",
			Unit: payload.NutrientUnit{
				Name: "g",
			},
		},
		NutrientValue: payload.NutrientValue{
			Value: float64(food.Zusatzinformationen.NwfettsaeurenDecimal1),
		},
	})

	nutrients = append(nutrients, payload.Nutrient{
		NutrientLabel: payload.NutrientLabel{
			Name: "Kohlenhydrate",
			Unit: payload.NutrientUnit{
				Name: "g",
			},
		},
		NutrientValue: payload.NutrientValue{
			Value: float64(food.Zusatzinformationen.NwkohlehydrateDecimal1),
		},
	})

	nutrients = append(nutrients, payload.Nutrient{
		NutrientLabel: payload.NutrientLabel{
			Name: "Zucker",
			Unit: payload.NutrientUnit{
				Name: "g",
			},
		},
		NutrientValue: payload.NutrientValue{
			Value: float64(food.Zusatzinformationen.NwzuckerDecimal1),
		},
	})

	nutrients = append(nutrients, payload.Nutrient{
		NutrientLabel: payload.NutrientLabel{
			Name: "Eiweiß",
			Unit: payload.NutrientUnit{
				Name: "g",
			},
		},
		NutrientValue: payload.NutrientValue{
			Value: float64(food.Zusatzinformationen.NweiweissDecimal1),
		},
	})

	nutrients = append(nutrients, payload.Nutrient{
		NutrientLabel: payload.NutrientLabel{
			Name: "Salz",
			Unit: payload.NutrientUnit{
				Name: "g",
			},
		},
		NutrientValue: payload.NutrientValue{
			Value: float64(food.Zusatzinformationen.NwsalzDecimal1),
		},
	})

	return &nutrients, nil
}


type LocalizedString struct {
	ValueDE *string
	ValueEN *string
}

type AdditiveResponse struct {
	ID int64 `json:"zusatzstoffeID"`
	Name string `json:"name"`
}

func ParseAdditives() (map[int64]LocalizedString, error) {
	additivesENResponse, err := sendRequestToSWT(AdditivesModel, NeuesPalais, EN)
	if err != nil {
		return nil, err
	}

	var additivesEN SWTResponse[AdditiveResponse]
	err = json.Unmarshal(additivesENResponse, &additivesEN)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	additivesDEResponse, err := sendRequestToSWT(AdditivesModel, NeuesPalais, DE)
	if err != nil {
		return nil, err
	}

	var additivesDE SWTResponse[AdditiveResponse]
	err = json.Unmarshal(additivesDEResponse, &additivesDE)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	additives := make(map[int64]LocalizedString)

	for _, add := range additivesEN.Content {
		additives[add.ID] = LocalizedString{
			ValueEN: &add.Name,
		}
	}

	for _, add := range additivesDE.Content {
		additives[add.ID] = LocalizedString{
			ValueEN: additives[add.ID].ValueEN,
			ValueDE: &add.Name,
		}
	}

	return additives, nil
}

type FeatureResponse struct {
	ID int64 `json:"gerichtmerkmalID"`
	Name string `json:"name"`
}

func ParseFeatures() (map[int64]LocalizedString, error) {
	featuresENResponse, err := sendRequestToSWT(FeaturesModel, NeuesPalais, EN)
	if err != nil {
		return nil, err
	}
	featuresDEResponse, err := sendRequestToSWT(FeaturesModel, NeuesPalais, DE)
	if err != nil {
		return nil, err
	}

	var featuresEN SWTResponse[FeatureResponse]
	err = json.Unmarshal(featuresENResponse, &featuresEN)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	var featuresDE SWTResponse[FeatureResponse]
	err = json.Unmarshal(featuresDEResponse, &featuresDE)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	features := make(map[int64]LocalizedString)

	for _, feature := range featuresEN.Content {
		features[feature.ID] = LocalizedString{
			ValueEN: &feature.Name,
		}
	}

	for _, feature := range featuresDE.Content {
		features[feature.ID] = LocalizedString{
			ValueEN: features[feature.ID].ValueEN,
			ValueDE: &feature.Name,
		}
	}

	return features, nil
}

func ExtractAdditives(food SpeiseplanGerichtDatum, additives map[int64]string) (*[]payload.Additive, error) {
	if len(*food.AdditivesIDsString) == 0 {
		return nil, nil
	}
	additivesArray := strings.Split(*food.AdditivesIDsString, ",")

	var result []payload.Additive

	for _, additiveID := range additivesArray {
		additiveIDInt, err := strconv.Atoi(additiveID)
		if err != nil {
			return nil, fmt.Errorf("error parsing additive ID: %w", err)
		}
		newAdditive := additives[int64(additiveIDInt)]
		result = append(result, payload.Additive{
			Name: newAdditive,
		})
	}

	return &result, nil
}
