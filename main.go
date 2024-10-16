package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"mensa-queue/internal/images"
	parsers "mensa-queue/internal/parse"
	"mensa-queue/internal/payload"
	"mensa-queue/internal/repository"
	"os"
	"time"
)

var recipes images.Recipes

func main() {
	ctx := context.Background()
	err := godotenv.Load() // 👈 load .env file
	if err != nil {
		log.Fatal(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")

	for {
		// Database connection
		conn, err := pgx.Connect(ctx, databaseURL)
		if err != nil {
			fmt.Printf("Unable to connect to database: %v\n", err)
			panic(err)
		}

		getAllMensas(ctx, conn)

		// TODO: Check if ComfyUI is reachable (only when my PC is on)
		go images.GenerateImages(recipes, ctx)

		conn.Close(ctx)

		time.Sleep(time.Hour)
	}
}

func getAllMensas(ctx context.Context, conn *pgx.Conn) {
	mensas := []payload.Mensa{payload.NeuesPalais}
	//mensas := []payload.Mensa{payload.NeuesPalais, payload.Griebnitzsee, payload.Golm, payload.Filmuniversitaet, payload.FHP, payload.Wildau, payload.Brandenburg}
	for _, mensa := range mensas {
		getMensaData(mensa, ctx, conn)
	}
}

func getMensaData(mensa payload.Mensa, ctx context.Context, conn *pgx.Conn) {
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
				Locales: []repository.InsertLocaleParams{
					{
						Name:   food.SpeiseplanAdvancedGericht.RecipeName,
						Locale: repository.EnumLocaleLocaleDe,
					},
					{
						Name:   food.Zusatzinformationen.GerichtnameAlternative,
						Locale: repository.EnumLocaleLocaleEn,
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

			recipeId, err := payload.InsertRecipe(recipe, t, languages, mensa, ctx, conn)
			if err != nil {
				fmt.Println("Error inserting recipe:", err)
				panic(err)
				continue
			}

			recipes = append(recipes, &images.RecipeData{
				ID:     recipeId,
				Prompt: food.Zusatzinformationen.GerichtnameAlternative,
			})
		}
	}
}
