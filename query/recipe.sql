-- name: FindAllRecipes :many
SELECT * FROM recipes;

-- name: InsertRecipe :one
INSERT INTO recipes (price_students, price_employees, price_guests, mensa_provider_id)
VALUES (@price_students::float8, @price_employees::float8, @price_guests::float8, $1)
RETURNING id;