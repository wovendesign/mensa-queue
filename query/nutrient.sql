-- name: FindOrCreateNutrientUnit :one
-- ON CONFLICT against nutrient_units_name_5_idx. The no-op
-- DO UPDATE makes the statement return the existing id on conflict.
INSERT INTO nutrient_units (name)
VALUES (@name)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING id;

-- name: FindOrCreateNutrientValue :one
-- ON CONFLICT against nutrient_values_value_5_idx. Values are deduped
-- globally, so identical numbers across dishes share one row.
INSERT INTO nutrient_values (value)
VALUES (@value::float8)
ON CONFLICT (value) DO UPDATE SET value = EXCLUDED.value
RETURNING id;

-- name: FindNutrientLabelByDeName :one
-- A nutrient label is identified by its localized name; the unique index
-- nutrient_labels_name_5_idx lives on (name, _locale).
SELECT _parent_id AS id
FROM nutrient_labels_locales
WHERE name = @name AND _locale = 'de'
LIMIT 1;

-- name: InsertNutrientLabel :one
INSERT INTO nutrient_labels (unit_id)
VALUES (@unit_id)
RETURNING id;

-- name: InsertNutrientLabelLocale :exec
INSERT INTO nutrient_labels_locales (name, _locale, _parent_id)
VALUES (@name, @locale::public._locales, @parent_id);

-- name: FindNutrientsForRecipe :many
SELECT nutrient_label_id, nutrient_value_id
FROM nutrients
WHERE recipe_id = @recipe_id;

-- name: DeleteNutrientsForRecipe :exec
DELETE FROM nutrients WHERE recipe_id = @recipe_id;

-- name: InsertNutrientForRecipe :exec
INSERT INTO nutrients (recipe_id, nutrient_value_id, nutrient_label_id)
VALUES (@recipe_id, @nutrient_value_id, @nutrient_label_id);
