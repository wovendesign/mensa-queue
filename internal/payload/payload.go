package payload

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"mensa-queue/internal/repository"
	"mensa-queue/models"
)

func InsertRecipe(recipe *models.Recipe, ctx context.Context, conn *pgx.Conn) (id *int32, err error) {
	repo := repository.New(conn)

	// Check if the recipe has "names"
	if len(recipe.Localization.Locales) == 0 {
		return nil, fmt.Errorf("no locale in recipe")
	}

	// The upstream dish ID is the primary identity of a recipe and survives
	// renames of the German dish name.
	var recipeID *int32
	if recipe.SourceID != "" {
		sourceRecipeID, err := repo.FindRecipeBySourceID(ctx, recipe.SourceID)
		switch {
		case err == nil:
			recipeID = &sourceRecipeID
		case errors.Is(err, pgx.ErrNoRows):
			// No recipe with this source ID yet -> fall back to the name lookup.
		default:
			return nil, fmt.Errorf("unable to find recipe by source ID: %w", err)
		}
	}

	// A found recipe may have been renamed upstream. Rename the linked German
	// locale in place before inserting locales, so insertLocales reuses that
	// row instead of creating a second German locale for the same recipe.
	if recipeID != nil {
		linkedLocales, err := repo.GetRecipeLocaleIDs(ctx, *recipeID)
		if err != nil {
			return nil, fmt.Errorf("unable to get recipe locales: %w", err)
		}
		newGermanName := recipe.Localization.Locales[0].Name
		for _, linkedLocale := range linkedLocales {
			if linkedLocale.Locale != repository.EnumLocaleLocaleDe {
				continue
			}
			if linkedLocale.Name != newGermanName {
				if err := repo.UpdateLocaleName(ctx, repository.UpdateLocaleNameParams{
					ID:   linkedLocale.ID,
					Name: newGermanName,
				}); err != nil {
					return nil, fmt.Errorf("unable to update recipe locale name: %w", err)
				}
			}
			break
		}
	}

	localeIDs, err := insertLocales(recipe.Localization.Locales, repo, ctx)
	if err != nil {
		return nil, err
	}

	// Legacy recipes were written before source IDs existed: find them by the
	// German name's linked locale and backfill the source ID.
	if recipeID == nil {
		recipeID, _ = repo.FindRecipeByLocale(ctx, localeIDs[0])
		if recipeID != nil && recipe.SourceID != "" {
			if err := repo.SetRecipeSourceID(ctx, repository.SetRecipeSourceIDParams{
				ID:       *recipeID,
				SourceID: recipe.SourceID,
			}); err != nil {
				return nil, fmt.Errorf("unable to set recipe source ID: %w", err)
			}
		}
	}

	if recipeID == nil {
		// No ID of a linked recipe was found -> create new recipe
		_recipeID, err := repo.InsertRecipe(ctx, repository.InsertRecipeParams{
			PriceStudents:   recipe.PriceStudents,
			PriceEmployees:  recipe.PriceEmployees,
			PriceGuests:     recipe.PriceGuests,
			MensaProviderID: *recipe.MensaProviderID,
			Category:        "main",
			SourceID:        recipe.SourceID,
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
				return nil, fmt.Errorf("unable to insert locale<->recipe rel: %v\n", err)
			}
		}
	} else {
		fmt.Printf("recipeID: %v\n", *recipeID)
		// The Recipe already existed, updating the prices, in case something changed
		err = repo.UpdateRecipePrices(ctx, repository.UpdateRecipePricesParams{
			ID:             *recipeID,
			PriceStudents:  recipe.PriceStudents,
			PriceEmployees: recipe.PriceEmployees,
			PriceGuests:    recipe.PriceGuests,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to update recipe: %v\n", err)
		}
	}

	// Insert Additional Stuff like Nutrients, Allergens, Additives, etc
	features := recipe.Localization.Features
	for _, feature := range features {
		localeIDs, err := insertLocales(feature, repo, ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to insert locale: %v\n", err)
		}

		// Check if Features exist already
		var featureID int32
		locale, _ := repo.FindLocale(ctx, feature[0].Name)
		if locale.FeaturesID == nil {
			_featureID, err := repo.InsertFeature(ctx)
			if err != nil {
				return nil, fmt.Errorf("unable to insert feature: %v\n", err)
			}
			featureID = _featureID
			//	Insert Locale Rels
			for _, localeID := range localeIDs {
				err = repo.InsertLocaleRel(ctx, repository.InsertLocaleRelParams{
					ParentID:  localeID,
					Path:      "feature",
					RecipeID:  nil,
					FeatureID: &featureID,
				})
				if err != nil {
					return nil, fmt.Errorf("unable to insert locale<->feature rel on feature: %v: %v\n", featureID, err)
				}
			}
		} else {
			featureID = *locale.FeaturesID
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
	_, err = insertServing(repo, ctx, recipe.Serving, recipeID)
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

func insertServing(repo *repository.Queries, ctx context.Context, serving *models.Serving, recipeID *int32) (*int32, error) {

	id, err := repo.InsertOrGetServing(ctx, repository.InsertOrGetServingParams{
		RecipeID: *recipeID,
		Date:     serving.Date,
		MensaID:  serving.MensaID,
	})
	if err != nil {
		fmt.Printf("Unable to insert serving: %v\n", err)
		return nil, err
	}
	return &id, nil
}
