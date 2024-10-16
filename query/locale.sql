-- name: FindLocale :one
SELECT
    locale.id, locale.name, locale.locale,
    locale_rels.path, locale_rels.recipes_id, locale_rels.features_id
FROM locale
         INNER JOIN locale_rels
                    ON locale.id = locale_rels.parent_id
WHERE locale.name = $1 LIMIT 1;

-- name: InsertLocale :one
INSERT INTO locale (name, locale)
VALUES ($1, $2)
RETURNING id;

-- name: InsertLocaleRel :exec
INSERT INTO locale_rels (parent_id, path, recipes_id, features_id)
VALUES ($1, $2, sqlc.narg(recipe_id), sqlc.narg(feature_id)::int);