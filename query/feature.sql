-- name: InsertFeature :one
INSERT INTO features (visible_small)
VALUES (false)
RETURNING id;

-- name: AddFeatureToRecipe :exec
INSERT INTO recipes_rels (parent_id, features_id, path)
VALUES ($1, $2, 'feature');