-- name: InsertServing :execrows
INSERT INTO servings (recipe_id, date, mensa_id)
SELECT $1, @date::timestamptz, @mensa_id
    WHERE NOT EXISTS (
    SELECT 1 FROM servings
    WHERE recipe_id = $1 AND date = @date::timestamptz AND mensa_id = @mensa_id
);

-- name: FindServing :one
SELECT id
FROM servings
WHERE recipe_id = $1 AND date = @date::timestamptz AND mensa_id = @mensa_id
LIMIT 1;