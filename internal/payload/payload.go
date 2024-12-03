package payload

import (
	"context"
	"fmt"
	"mensa-queue/internal/repository"
	"time"

	"github.com/jackc/pgx/v5"
)

type LocalNutrient struct {
	Unit    string
	Value   float64
	Locales []*repository.InsertLocaleParams
}

type LocalRecipe struct {
	Locales   []*repository.InsertLocaleParams
	Allergen  [][]*repository.InsertLocaleParams
	Additives [][]*repository.InsertLocaleParams
	Features  [][]*repository.InsertLocaleParams
	Nutrients []*LocalNutrient
	Recipe    Recipe
}

func InsertRecipe(recipe *LocalRecipe, date time.Time, mensa Mensa, ctx context.Context, conn *pgx.Conn) (id *int32, err error) {
	repo := repository.New(conn)

	// Check if the recipe has "names"
	if len(recipe.Locales) == 0 {
		return nil, fmt.Errorf("no locale in recipe")
	}

	localeIDs, err := insertLocales(recipe.Locales, repo, ctx)
	if err != nil {
		return nil, err
	}

	// Check if the locale has a Recipe linked to it
	// This is e.g. not the case if the Locales were newly created
	recipeID, _ := repo.FindRecipeByLocale(ctx, localeIDs[0])
	//if err != nil {
	//	fmt.Printf("error finding recipe: %v\n", err)
	//	//err = nil
	//	return nil, err
	//}

	fmt.Printf("recipeID: %v\n", recipeID)

	if recipeID == nil {
		// No ID of a linked recipe was found -> create new recipe
		_recipeID, err := repo.InsertRecipe(ctx, repository.InsertRecipeParams{
			PriceStudents:   recipe.Recipe.PriceStudents,
			PriceEmployees:  recipe.Recipe.PriceEmployees,
			PriceGuests:     recipe.Recipe.PriceGuests,
			MensaProviderID: 1,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to insert recipe: %v\n", err)
		}
		recipeID = &_recipeID

		// Link the Recipe to the Locales
		for _, localeID := range localeIDs {
			err = repo.InsertLocaleRel(ctx, repository.InsertLocaleRelParams{
				ParentID:  localeID,
				Path:      "recipe",
				RecipeID:  recipeID,
				FeatureID: nil,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to insert locale: %v\n", err)
			}
		}
	} else {
		// The Recipe already existed, updating the prices, in case something changed
		err = repo.UpdateRecipePrices(ctx, repository.UpdateRecipePricesParams{
			ID:             *recipeID,
			PriceStudents:  recipe.Recipe.PriceStudents,
			PriceEmployees: recipe.Recipe.PriceEmployees,
			PriceGuests:    recipe.Recipe.PriceGuests,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to update recipe: %v\n", err)
		}
	}

	// Insert Additional Stuff like Nutrients, Allergens, Additives, etc
	features := recipe.Features
	for _, feature := range features {
		localeIDs, err := insertLocales(feature, repo, ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to insert locale: %v\n", err)
		}

		// Check if Features exist already
		locale, _ := repo.FindLocale(ctx, feature[0].Name)
		if locale.FeaturesID != nil {
			continue
		}

		featureID, err := repo.InsertFeature(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to insert feature: %v\n", err)
		}

		//	Insert Locale Rels
		for _, localeID := range localeIDs {
			err = repo.InsertLocaleRel(ctx, repository.InsertLocaleRelParams{
				ParentID:  localeID,
				Path:      "feature",
				RecipeID:  nil,
				FeatureID: &featureID,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to insert locale: %v\n", err)
			}
		}

		//	Add Feature to Recipe
		err = repo.AddFeatureToRecipe(ctx, repository.AddFeatureToRecipeParams{
			ParentID:   *recipeID,
			FeaturesID: &featureID,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to insert feature<->recipe rel: %v\n", err)
		}
	}

	// Create a Serving of the Recipe
	// - A recipe has no dates or mensas attached to it
	// This creates a new serving, if there wasn't one already
	_, err = insertServing(mensa, err, repo, ctx, recipeID, date)
	if err != nil {
		return nil, err
	}

	return recipeID, nil
}

func insertLocales(locales []*repository.InsertLocaleParams, repo *repository.Queries, ctx context.Context) ([]int32, error) {
	// This array will be populated with the IDs in the database of the Locales
	localeIDs := make([]int32, 0)

	// Loop over all available locales
	for _, locale := range locales {
		// Insert the locale if it doesn't exist, return it if it does
		id, err := repo.InsertLocaleIfNotExists(ctx, repository.InsertLocaleIfNotExistsParams{
			Name:   locale.Name,
			Locale: locale.Locale,
		})
		if err != nil {
			fmt.Printf("error inserting locale: %v\n", err)
			return nil, err
		}

		// Append the ID to the ID Array
		localeIDs = append(localeIDs, id)
	}
	return localeIDs, nil
}

func insertServing(mensa Mensa, err error, repo *repository.Queries, ctx context.Context, recipeID *int32, date time.Time) (*int32, error) {
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

	id, err := repo.InsertOrGetServing(ctx, repository.InsertOrGetServingParams{
		RecipeID: *recipeID,
		Date:     date,
		MensaID:  &mensaId,
	})
	if err != nil {
		fmt.Printf("Unable to insert serving: %v\n", err)
		return nil, err
	}
	return &id, nil
}
