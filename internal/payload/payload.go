package payload

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"mensa-queue/internal/repository"
	"time"
)

func InsertRecipe(recipe *LocalRecipe, date time.Time, language []Language, mensa Mensa, ctx context.Context, conn *pgx.Conn) (id *int, err error) {
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
		locales = append(locales, _locale)
	}

	var recipeID int32

	fmt.Println(len(locales))

	if len(locales) == 0 {

		// Create New Recipe
		recipeID, err = repo.InsertRecipe(ctx, repository.InsertRecipeParams{
			PriceStudents:   pgtype.Numeric{},
			PriceEmployees:  pgtype.Numeric{},
			PriceGuests:     pgtype.Numeric{},
			MensaProviderID: 2,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to insert recipe: %v\n", err)
		}
	} else {
		if locales[0].RecipesID.Valid {
			recipeID = locales[0].RecipesID.Int32
		} else {
			return nil, fmt.Errorf("invalid recipe id")
		}
		// Recipe Already Exists
		// Create Serving if it doesn't exist
	}

	rows, err := repo.InsertServing(ctx, repository.InsertServingParams{
		RecipeID: recipeID,
		Date: pgtype.Timestamptz{
			Time: date,
		},
		MensaID: pgtype.Int4{Int32: int32(mensa)},
	})
	if err != nil {
		fmt.Printf("Unable to insert recipe: %v\n", err)
		return nil, err
	}
	if rows == 0 {
		return nil, fmt.Errorf("no rows inserted")
	}

	return nil, nil
}
