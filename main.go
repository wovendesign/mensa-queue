package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"mensa-queue/internal/config"
	"mensa-queue/internal/images"
	parsers "mensa-queue/internal/parse"
	"mensa-queue/internal/payload"
	"mensa-queue/internal/repository"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var recipes images.Recipes

func loadConfig() (*pgx.ConnConfig, error) {
	cfg, err := config.NewDatabase()
	if err != nil {
		return nil, err
	}

	//conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	return pgx.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	))
}

func main() {
	ctx := context.Background()

	err := godotenv.Load() // 👈 load .env file
	if err != nil {
		log.Printf("Error loading .env file\n")
	}

	for {
		// Database connection
		//pgConfig, err := loadConfig()
		//if err != nil {
		//	log.Fatal(err)
		//}
		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		//conn, err := pgx.ConnectConfig(ctx, pgConfig)
		if err != nil {
			fmt.Printf("Unable to connect to database: %v\n", err)
			panic(err)
		}

		getAllMensas(ctx, conn)

		// TODO: Check if ComfyUI is reachable (only when my PC is on)
		// TODO: Check if AI Image already exists before generating one
		// go images.GenerateImages(recipes, ctx)

		conn.Close(ctx)

		time.Sleep(time.Hour)
	}
}

func getAllMensas(ctx context.Context, conn *pgx.Conn) {
	//mensas := []payload.Mensa{payload.NeuesPalais}
	mensas := []payload.Mensa{payload.NeuesPalais, payload.Griebnitzsee, payload.Golm, payload.Filmuniversitaet, payload.FHP, payload.Wildau, payload.Brandenburg}
	for _, mensa := range mensas {
		getMensaData(mensa, ctx, conn)
	}
}

func getMensaData(mensa payload.Mensa, ctx context.Context, conn *pgx.Conn) {
	languages := []repository.EnumLocaleLocale{repository.EnumLocaleLocaleDe, repository.EnumLocaleLocaleEn}
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

	categoryMap, err := parsers.ParseMealCategory(mensa)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, week := range *foodContent {
		for _, food := range week.SpeiseplanGerichtData {
			if food.Zusatzinformationen.MitarbeiterpreisDecimal2 == 0 || strings.Contains(food.SpeiseplanAdvancedGericht.RecipeName, "Preis pro") {
				continue
			}

			category, err := parsers.ExtractCategories(food, categoryMap)
			if err != nil {
				log.Fatal(err)
				return
			}

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
				Recipe: payload.Recipe{
					PriceStudents: &food.Zusatzinformationen.MitarbeiterpreisDecimal2,
					PriceGuests:   &food.Zusatzinformationen.GaestepreisDecimal2,
					MensaProvider: 1,
				},
				Nutrients: nutrients,
				Allergen:  allergens,
				Additives: additives,
				Features:  features,
				Category:  category,
			}

			t, err := time.Parse(time.RFC3339, food.SpeiseplanAdvancedGericht.Date)
			t = t.UTC()
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}

			fmt.Printf("Recipe: %v\n", recipe)

			recipeId, err := payload.InsertRecipe(recipe, t, mensa, ctx, conn)
			if err != nil {
				fmt.Println("Error inserting recipe:", err)
				continue
			}

			recipes = append(recipes, &images.RecipeData{
				ID:     recipeId,
				Prompt: food.Zusatzinformationen.GerichtnameAlternative,
			})
		}
	}
}
