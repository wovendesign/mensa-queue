-- name: InsertFeature :one
INSERT INTO features (visible_small)
VALUES (false)
RETURNING id;

-- name: AddFeatureToRecipe :exec
INSERT INTO recipes_rels (parent_id, features_id, path)
SELECT $1, $2, 'feature'
WHERE NOT EXISTS (
    SELECT 1 FROM recipes_rels WHERE parent_id = $1 AND features_id = $2
);