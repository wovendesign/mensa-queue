-- name: InsertMensaProvider :one
WITH ins AS (
    INSERT INTO mensa_provider (name, description, slug, uuid)
        SELECT $1, $2, $3, @uuid::text
        WHERE NOT EXISTS (
            SELECT 1 FROM mensa_provider
            WHERE uuid = @uuid::text
        )
        RETURNING id
)
SELECT id FROM ins
UNION
SELECT id FROM mensa_provider WHERE uuid = @uuid::text;