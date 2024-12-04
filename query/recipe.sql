-- name: FindAllRecipes :many
SELECT * FROM recipes;

-- name: InsertRecipe :one
INSERT INTO recipes (price_students, price_employees, price_guests, mensa_provider_id, category)
VALUES (sqlc.narg(price_students)::float8, sqlc.narg(price_employees)::float8, sqlc.narg(price_guests)::float8, $1, $2)
RETURNING id;

-- name: UpdateRecipePrices :exec
UPDATE recipes
SET price_students = sqlc.narg(price_students)::float8, price_employees = sqlc.narg(price_employees)::float8, price_guests = sqlc.narg(price_guests)::float8
WHERE id = $1;

-- name: SetRecipeAIImage :exec
UPDATE recipes
SET ai_thumbnail_id = @ai_thumbnail_id::int
WHERE id = $1;