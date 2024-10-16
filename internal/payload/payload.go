package payload

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"mensa-queue/internal/repository"
	"time"
)

func InsertRecipe(recipe *LocalRecipe, date time.Time, language []Language, mensa Mensa, ctx context.Context, conn *pgx.Conn) (id *int32, err error) {
	repo := repository.New(conn)

	if len(recipe.Locales) == 0 {
		return nil, fmt.Errorf("no locale in recipe")
	}

	locales := make([]repository.FindLocaleRow, 0)

	for _, locale := range recipe.Locales {
		_locale, err := repo.FindLocale(ctx, locale.Name)
		if err != nil {
			fmt.Printf("Unable to find locale: %v\n", err)
			continue
		}
		fmt.Printf("Found locale: %+v\n", _locale)
		locales = append(locales, _locale)
	}

	var recipeID int32

	if len(locales) == 0 {
		// Create New Recipe
		recipeID, err = repo.InsertRecipe(ctx, repository.InsertRecipeParams{
			PriceStudents:   recipe.Recipe.PriceStudents,
			PriceEmployees:  recipe.Recipe.PriceEmployees,
			PriceGuests:     recipe.Recipe.PriceGuests,
			MensaProviderID: 1,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to insert recipe: %v\n", err)
		}

		// Insert Locales
		for _, locale := range recipe.Locales {

			id2, err := repo.InsertLocale(ctx, repository.InsertLocaleParams{
				Name:   locale.Name,
				Locale: locale.Locale,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to insert locale: %v\n", err)
			}

			err = repo.InsertLocaleRel(ctx, repository.InsertLocaleRelParams{
				ParentID:  id2,
				Path:      "recipe",
				RecipeID:  &recipeID,
				FeatureID: nil,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to insert locale: %v\n", err)
			}
		}
	} else {

		if locales[0].RecipesID != nil {
			recipeID = *locales[0].RecipesID
		} else {
			return nil, fmt.Errorf("invalid recipe id")
		}

		err = repo.UpdateRecipePrices(ctx, repository.UpdateRecipePricesParams{
			ID:             recipeID,
			PriceStudents:  recipe.Recipe.PriceStudents,
			PriceEmployees: recipe.Recipe.PriceEmployees,
			PriceGuests:    recipe.Recipe.PriceGuests,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to update recipe: %v\n", err)
		}
		// Recipe Already Exists
		// Create Serving if it doesn't exist
		fmt.Println(recipeID)
	}

	mensaMap := map[Mensa]int32{
		NeuesPalais:      1,
		Griebnitzsee:     2,
		Golm:             3,
		Filmuniversitaet: 4,
		FHP:              5,
		Wildau:           6,
		Brandenburg:      7,
	}
	mensaId := mensaMap[mensa]

	rows, err := repo.InsertServing(ctx, repository.InsertServingParams{
		RecipeID: recipeID,
		Date:     date,
		MensaID:  &mensaId,
	})
	if err != nil {
		fmt.Printf("Unable to insert recipe: %v\n", err)
		return nil, err
	}
	if rows == 0 {
		// Check if Serving already exists, if not, return error
		serving, err := repo.FindServing(ctx, repository.FindServingParams{
			RecipeID: recipeID,
			Date:     date,
			MensaID:  &mensaId,
		})
		if err != nil {
			fmt.Printf("Unable to find serving: %v\n", err)
		}
		fmt.Printf("Found serving: %+v\n", serving)
		//return nil, fmt.Errorf("no rows inserted: %+v\n", err)
	}

	return &recipeID, nil
}
