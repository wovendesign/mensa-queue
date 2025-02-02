-- name: InsertMensa :one
WITH ins AS (
    INSERT INTO mensa (name, description, slug, uuid, address_latitude, address_longitude, provider_id)
        SELECT $1, $2, $3, @uuid::text, 0, 0, $4
        WHERE NOT EXISTS (
            SELECT 1 FROM mensa
            WHERE uuid = @uuid::text
        )
        RETURNING id
)
SELECT id FROM ins
UNION
SELECT id FROM mensa WHERE uuid = @uuid::text AND provider_id = $4;