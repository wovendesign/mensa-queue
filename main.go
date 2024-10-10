package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"mensa-queue/internal/images"
	parsers "mensa-queue/internal/parse"
	"mensa-queue/internal/payload"

	"time"

	"gorm.io/gorm"
)

var recipes images.Recipes

func main() {
	for {
		// Database connection
		dsn := "host=127.0.0.1 user=mensauser password=postgres dbname=mensahhub port=5432 sslmode=disable TimeZone=Europe/Berlin"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		getAllMensas(db)

		// TODO: Check if ComfyUI is reachable (only when my PC is on)
		go images.GenerateImages(recipes)

		if conn, err := db.DB(); err == nil {
			_ = conn.Close()
		}

		time.Sleep(time.Hour)
	}
}

func getAllMensas(db *gorm.DB) {
	mensas := []payload.Mensa{payload.NeuesPalais, payload.Griebnitzsee, payload.Golm, payload.Filmuniversitaet, payload.FHP, payload.Wildau, payload.Brandenburg}
	for _, mensa := range mensas {
		getMensaData(mensa, db)
	}
}

func getMensaData(mensa payload.Mensa, db *gorm.DB) {
	languages := []payload.Language{payload.EN, payload.DE}
	foodContent, err := parsers.ParsePotsdamMensaData(mensa)
	if err != nil {
		log.Fatal(err)
	}

	additiveMap, err := parsers.ParseAdditives(languages, mensa)
	if err != nil {
		log.Fatal(err)
		return
	}

	allergenMap, err := parsers.ParseAllergens(languages, mensa)
	if err != nil {
		log.Fatal(err)
		return
	}
	// fmt.Printf("%+v\n", allergenMap
	featureMap, err := parsers.ParseFeatures(languages, mensa)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, week := range *foodContent {
		for _, food := range week.SpeiseplanGerichtData {
			// fmt.Printf("%+v\n", food)

			nutrients, err := parsers.ExtractNutrients(food)
			if err != nil {
				log.Fatal(err)
				return
			}

			additives, err := parsers.ExtractAdditives(food, additiveMap, languages)
			if err != nil {
				log.Fatal(err)
				return
			}

			allergens, err := parsers.ExtractAllergens(food, allergenMap, languages)
			if err != nil {
				log.Fatal(err)
				return
			}

			features, err := parsers.ExtractFeatures(food, featureMap, languages)
			if err != nil {
				log.Fatal(err)
				return
			}

			recipe := &payload.LocalRecipe{
				Locales: []payload.Locale{
					{
						Name:   food.SpeiseplanAdvancedGericht.RecipeName,
						Locale: "de",
					},
					{
						Name:   food.Zusatzinformationen.GerichtnameAlternative,
						Locale: "en",
					},
				},
				Recipe: payload.Recipe{
					PriceStudents: &food.Zusatzinformationen.MitarbeiterpreisDecimal2,
					PriceGuests:   &food.Zusatzinformationen.GaestepreisDecimal2,
					MensaProvider: 1,
				},
				Nutrients: nutrients,
				Allergen:  allergens,
				Additives: additives,
				Features:  features,
			}

			t, err := time.Parse(time.RFC3339, food.SpeiseplanAdvancedGericht.Date)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}

			payload.InsertRecipe(recipe, t, languages, mensa, db)

			recipes = append(recipes, &images.RecipeData{
				ID:     recipe.Recipe.ID,
				Prompt: food.Zusatzinformationen.GerichtnameAlternative,
			})
		}
	}
}
