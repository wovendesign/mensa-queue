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

-- name: InsertLocaleIfNotExists :one
WITH ins AS (
    INSERT INTO locale (name, locale)
        SELECT sqlc.arg(name)::varchar, sqlc.arg(locale)
        WHERE NOT EXISTS (
            SELECT 1 FROM locale WHERE name = sqlc.arg(name)::varchar AND locale = sqlc.arg(locale)
        )
        RETURNING id
)
SELECT id FROM ins
UNION
SELECT id FROM locale WHERE name = sqlc.arg(name)::varchar AND locale = sqlc.arg(locale);

-- name: InsertLocaleRel :exec
INSERT INTO locale_rels (parent_id, path, recipes_id, features_id)
VALUES ($1, $2, sqlc.narg(recipe_id), sqlc.narg(feature_id)::int);

-- name: FindRecipeByLocale :one
SELECT locale_rels.recipes_id
from locale_rels
WHERE locale_rels.parent_id = $1 AND locale_rels.path = 'recipe'
LIMIT 1;