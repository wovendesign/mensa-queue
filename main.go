package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"mensa-queue/adapters"
	stw_brandenburg_west "mensa-queue/adapters/stw-brandenburg-west"
	"mensa-queue/internal/config"
	"mensa-queue/internal/images"
	"mensa-queue/internal/payload"
	"mensa-queue/internal/repository"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var recipes images.Recipes

func loadConfig() (*pgx.ConnConfig, error) {
	cfg, err := config.NewDatabase()
	if err != nil {
		return nil, err
	}

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

	providerAdapters := []adapters.Adapter{
		stw_brandenburg_west.NewAdapter("Studierendenwerk Brandenburg West"),
	}

	// Database connection
	//pgConfig, err := loadConfig()
	//if err != nil {
	//	log.Fatal(err)
	//}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		panic(err)
	}

	for _, providerAdapter := range providerAdapters {
		err = providerAdapter.RegisterAdapter(ctx, conn)
		if err != nil {
			fmt.Printf("Unable to register adapter: %v\n", err)
			continue
		}

		mensas := providerAdapter.GetAllMensas()

		for _, mensa := range mensas {
			err = mensa.RegisterMensa(ctx, conn)
			if err != nil {
				fmt.Printf("Unable to register mensa: %v\n", err)
			}
		}
	}

	conn.Close(ctx)

	for {
		// Database connection
		pgConfig, err := loadConfig()
		if err != nil {
			log.Fatal(err)
		}
		conn, err := pgx.ConnectConfig(ctx, pgConfig)
		if err != nil {
			fmt.Printf("Unable to connect to database: %v\n", err)
			panic(err)
		}

		for _, providerAdapter := range providerAdapters {
			for _, mensa := range providerAdapter.GetAllMensas() {
				if !mensa.IsRegistered() {
					continue
				}
				menu, err := mensa.ParseMenu()
				if err != nil {
					fmt.Printf("Unable to parse menu: %v\n", err)
				}

				for _, recipe := range menu {
					recipeId, err := payload.InsertRecipe(recipe, ctx, conn)
					if err != nil {
						fmt.Println("Error inserting recipe:", err)
						continue
					}

					if mensa.AiGenerationEnabled() {
						var enRecipeName string

						for _, locale := range recipe.Localization.Locales {
							if locale.Locale == repository.EnumLocaleLocaleEn {
								enRecipeName = locale.Name
							}
						}

						recipes = append(recipes, &images.RecipeData{
							ID:     recipeId,
							Prompt: enRecipeName,
						})
					}
				}
			}
		}

		// TODO: Check if ComfyUI is reachable (only when my PC is on)
		// TODO: Check if AI Image already exists before generating one
		// go images.GenerateImages(recipes, ctx)

		err = conn.Close(ctx)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Hour)
	}
}
