-- name: FindAllRecipes :many
SELECT * FROM recipes;

-- name: FindLocale :one
SELECT
    locale.id, locale.name, locale.locale,
    locale_rels.path, locale_rels.recipes_id, locale_rels.features_id
FROM locale
    INNER JOIN locale_rels
    ON locale.id = locale_rels.parent_id
WHERE locale.name = $1 LIMIT 1;

-- name: InsertServing :execrows
INSERT INTO servings (recipe_id, date, mensa_id)
SELECT $1, $2, $3
    WHERE NOT EXISTS (
    SELECT 1 FROM servings
    WHERE recipe_id = $1 AND date = $2 AND mensa_id = $3
);

-- name: InsertRecipe :one
INSERT INTO recipes (price_students, price_employees, price_guests, mensa_provider_id)
VALUES ($1, $2, $3, $4)
RETURNING id;