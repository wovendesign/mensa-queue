-- name: FindAllRecipes :many
SELECT * FROM recipes;

-- name: FindRecipeByID :one
SELECT * FROM recipes WHERE id = $1;

-- name: FindRecipeBySourceID :one
SELECT id FROM recipes WHERE source_id = sqlc.arg(source_id)::text LIMIT 1;

-- name: InsertRecipe :one
INSERT INTO recipes (price_students, price_employees, price_guests, mensa_provider_id, category, source_id)
VALUES (sqlc.narg(price_students)::float8, sqlc.narg(price_employees)::float8, sqlc.narg(price_guests)::float8, $1, $2, sqlc.arg(source_id)::text)
RETURNING id;

-- name: UpdateRecipePrices :exec
UPDATE recipes
SET price_students = sqlc.narg(price_students)::float8, price_employees = sqlc.narg(price_employees)::float8, price_guests = sqlc.narg(price_guests)::float8
WHERE id = $1;

-- name: SetRecipeSourceID :exec
UPDATE recipes
SET source_id = sqlc.arg(source_id)::text
WHERE id = $1;

-- name: GetRecipeLocaleIDs :many
SELECT locale.id, locale.name, locale.locale
FROM locale
         INNER JOIN locale_rels ON locale.id = locale_rels.parent_id
WHERE locale_rels.path = 'recipe' AND locale_rels.recipes_id = sqlc.arg(recipes_id)::int;

-- name: UpdateLocaleName :exec
UPDATE locale
SET name = $2
WHERE id = $1;

-- name: SetRecipeAIImage :exec
UPDATE recipes
SET ai_thumbnail_id = @ai_thumbnail_id::int
WHERE id = $1;
