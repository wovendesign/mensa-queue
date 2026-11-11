-- Allergens and additives are localized via their own *_locales tables
-- (name, _locale, _parent_id). Those tables have no unique index on the
-- name, so dedupe works by find-by-de-name-then-insert. Linking to the
-- recipe goes through recipes_rels with path = field name of the Recipe
-- collection relationship ('allergens' / 'additives').

-- name: FindAllergenByDeName :one
SELECT _parent_id AS id
FROM allergens_locales
WHERE name = @name AND _locale = 'de'
LIMIT 1;

-- name: InsertAllergen :one
INSERT INTO allergens DEFAULT VALUES
RETURNING id;

-- name: InsertAllergenLocale :exec
INSERT INTO allergens_locales (name, _locale, _parent_id)
VALUES (@name, @locale::public._locales, @parent_id);

-- name: AddAllergenToRecipe :exec
INSERT INTO recipes_rels (parent_id, path, allergens_id)
SELECT @recipe_id, 'allergens', @allergens_id
WHERE NOT EXISTS (
    SELECT 1 FROM recipes_rels
    WHERE parent_id = @recipe_id AND allergens_id = @allergens_id AND path = 'allergens'
);

-- name: FindAdditiveByDeName :one
SELECT _parent_id AS id
FROM additives_locales
WHERE name = @name AND _locale = 'de'
LIMIT 1;

-- name: InsertAdditive :one
INSERT INTO additives DEFAULT VALUES
RETURNING id;

-- name: InsertAdditiveLocale :exec
INSERT INTO additives_locales (name, _locale, _parent_id)
VALUES (@name, @locale::public._locales, @parent_id);

-- name: AddAdditiveToRecipe :exec
INSERT INTO recipes_rels (parent_id, path, additives_id)
SELECT @recipe_id, 'additives', @additives_id
WHERE NOT EXISTS (
    SELECT 1 FROM recipes_rels
    WHERE parent_id = @recipe_id AND additives_id = @additives_id AND path = 'additives'
);
