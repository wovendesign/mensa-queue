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
	err = persistNutrients(repo, conn, ctx, *recipeID, recipe.Localization.Nutrients)
	if err != nil {
		return nil, fmt.Errorf("unable to persist nutrients: %v\n", err)
	}

	err = persistAllergens(repo, ctx, *recipeID, recipe.Localization.Allergen)
	if err != nil {
		return nil, fmt.Errorf("unable to persist allergens: %v\n", err)
	}

	err = persistAdditives(repo, ctx, *recipeID, recipe.Localization.Additives)
	if err != nil {
		return nil, fmt.Errorf("unable to persist additives: %v\n", err)
	}

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

// persistNutrients replaces the nutrients of a recipe with the parsed values.
// SWP sends 0 (never null) for missing nutrition data, so only values > 0 are
// stored; if everything is 0, the recipe ends up without nutrient rows.
// Values, units and labels are global lookup tables and are deduped via their
// unique indexes, so identical data across dishes shares rows. The recipe's
// rows are only rewritten when the set of (label, value) pairs actually
// changed, to avoid churn on repeated runs.
func persistNutrients(repo *repository.Queries, conn *pgx.Conn, ctx context.Context, recipeID int32, nutrients []*models.NutrientLocalization) error {
	type nutrientRow struct {
		LabelID int32
		ValueID int32
	}

	desired := make([]nutrientRow, 0, len(nutrients))
	for _, nutrient := range nutrients {
		if nutrient.Value <= 0 {
			continue
		}

		labelID, err := findOrCreateNutrientLabel(repo, ctx, nutrient)
		if err != nil {
			return fmt.Errorf("unable to find or create nutrient label: %v\n", err)
		}

		valueID, err := repo.FindOrCreateNutrientValue(ctx, nutrient.Value)
		if err != nil {
			return fmt.Errorf("unable to find or create nutrient value: %v\n", err)
		}

		desired = append(desired, nutrientRow{LabelID: labelID, ValueID: valueID})
	}

	// Skip the rewrite if nothing changed
	existingRows, err := repo.FindNutrientsForRecipe(ctx, &recipeID)
	if err != nil {
		return fmt.Errorf("unable to find existing nutrients: %v\n", err)
	}
	existing := make(map[[2]int32]struct{}, len(existingRows))
	for _, row := range existingRows {
		if row.NutrientLabelID == nil || row.NutrientValueID == nil {
			continue
		}
		existing[[2]int32{*row.NutrientLabelID, *row.NutrientValueID}] = struct{}{}
	}
	if len(existing) == len(desired) {
		unchanged := true
		for _, row := range desired {
			if _, ok := existing[[2]int32{row.LabelID, row.ValueID}]; !ok {
				unchanged = false
				break
			}
		}
		if unchanged {
			return nil
		}
	}

	// Replace the recipe's nutrient rows atomically
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v\n", err)
	}
	defer tx.Rollback(ctx)

	txRepo := repository.New(tx)
	err = txRepo.DeleteNutrientsForRecipe(ctx, &recipeID)
	if err != nil {
		return fmt.Errorf("unable to delete old nutrients: %v\n", err)
	}

	for _, row := range desired {
		err = txRepo.InsertNutrientForRecipe(ctx, repository.InsertNutrientForRecipeParams{
			RecipeID:        &recipeID,
			NutrientValueID: &row.ValueID,
			NutrientLabelID: &row.LabelID,
		})
		if err != nil {
			return fmt.Errorf("unable to insert nutrient: %v\n", err)
		}
	}

	return tx.Commit(ctx)
}

// findOrCreateNutrientLabel resolves a nutrient label by its German name
// (the unique index nutrient_labels_name_5_idx lives on (name, _locale)) and
// creates the label with its localized names if it doesn't exist yet.
func findOrCreateNutrientLabel(repo *repository.Queries, ctx context.Context, nutrient *models.NutrientLocalization) (int32, error) {
	deName, ok := deName(nutrient.Locales)
	if !ok {
		return 0, fmt.Errorf("no german name for nutrient")
	}

	labelID, err := repo.FindNutrientLabelByDeName(ctx, deName)
	if err == nil {
		return labelID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	// Not found -> create the label together with its unit and localized names
	unitID, err := repo.FindOrCreateNutrientUnit(ctx, nutrient.Unit)
	if err != nil {
		return 0, fmt.Errorf("unable to find or create nutrient unit: %v\n", err)
	}

	labelID, err = repo.InsertNutrientLabel(ctx, unitID)
	if err != nil {
		return 0, fmt.Errorf("unable to insert nutrient label: %v\n", err)
	}

	for _, locale := range nutrient.Locales {
		err = repo.InsertNutrientLabelLocale(ctx, repository.InsertNutrientLabelLocaleParams{
			Name:     locale.Name,
			Locale:   repository.Locales(locale.Locale),
			ParentID: labelID,
		})
		if err != nil {
			return 0, fmt.Errorf("unable to insert nutrient label locale: %v\n", err)
		}
	}

	return labelID, nil
}

// persistAllergens dedupes allergens by their German name and links each one
// to the recipe idempotently via recipes_rels (path = 'allergens').
func persistAllergens(repo *repository.Queries, ctx context.Context, recipeID int32, allergens [][]*repository.InsertLocaleParams) error {
	for _, locales := range allergens {
		allergenID, err := findOrCreateLocalizedEntity(
			repo, ctx, locales,
			repo.FindAllergenByDeName,
			repo.InsertAllergen,
			func(ctx context.Context, entityID int32, locale *repository.InsertLocaleParams) error {
				return repo.InsertAllergenLocale(ctx, repository.InsertAllergenLocaleParams{
					Name:     &locale.Name,
					Locale:   repository.Locales(locale.Locale),
					ParentID: entityID,
				})
			},
		)
		if err != nil {
			return err
		}

		err = repo.AddAllergenToRecipe(ctx, repository.AddAllergenToRecipeParams{
			RecipeID:    recipeID,
			AllergensID: &allergenID,
		})
		if err != nil {
			return fmt.Errorf("unable to add allergen to recipe: %v\n", err)
		}
	}
	return nil
}

// persistAdditives is the additive twin of persistAllergens.
func persistAdditives(repo *repository.Queries, ctx context.Context, recipeID int32, additives [][]*repository.InsertLocaleParams) error {
	for _, locales := range additives {
		additiveID, err := findOrCreateLocalizedEntity(
			repo, ctx, locales,
			repo.FindAdditiveByDeName,
			repo.InsertAdditive,
			func(ctx context.Context, entityID int32, locale *repository.InsertLocaleParams) error {
				return repo.InsertAdditiveLocale(ctx, repository.InsertAdditiveLocaleParams{
					Name:     &locale.Name,
					Locale:   repository.Locales(locale.Locale),
					ParentID: entityID,
				})
			},
		)
		if err != nil {
			return err
		}

		err = repo.AddAdditiveToRecipe(ctx, repository.AddAdditiveToRecipeParams{
			RecipeID:    recipeID,
			AdditivesID: &additiveID,
		})
		if err != nil {
			return fmt.Errorf("unable to add additive to recipe: %v\n", err)
		}
	}
	return nil
}

// findOrCreateLocalizedEntity resolves an allergen/additive by its German
// name. Those lookup tables have no unique index on the name, so dedupe works
// by find-then-insert (a single worker runs against the DB, so there is no
// race to guard against).
func findOrCreateLocalizedEntity(
	repo *repository.Queries,
	ctx context.Context,
	locales []*repository.InsertLocaleParams,
	findByDeName func(context.Context, *string) (int32, error),
	insertEntity func(context.Context) (int32, error),
	insertLocale func(context.Context, int32, *repository.InsertLocaleParams) error,
) (int32, error) {
	de, ok := deName(locales)
	if !ok {
		return 0, fmt.Errorf("no german name for entity")
	}

	entityID, err := findByDeName(ctx, &de)
	if err == nil {
		return entityID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	entityID, err = insertEntity(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to insert entity: %v\n", err)
	}

	for _, locale := range locales {
		err = insertLocale(ctx, entityID, locale)
		if err != nil {
			return 0, fmt.Errorf("unable to insert entity locale: %v\n", err)
		}
	}

	return entityID, nil
}

// deName returns the German name of a set of localized names.
func deName(locales []*repository.InsertLocaleParams) (string, bool) {
	for _, locale := range locales {
		if locale.Locale == repository.EnumLocaleLocaleDe {
			return locale.Name, true
		}
	}
	return "", false
}
