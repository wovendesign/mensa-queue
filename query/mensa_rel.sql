-- name: GetMensaOpeningDayIDs :many
SELECT opening_day_id
FROM mensa_rels
WHERE parent_id = $1 AND path = 'opening_days' AND opening_day_id IS NOT NULL
ORDER BY "order", id;

-- name: DeleteMensaOpeningDayRels :exec
DELETE FROM mensa_rels
WHERE parent_id = $1 AND path = 'opening_days';

-- name: InsertMensaOpeningDayRel :exec
INSERT INTO mensa_rels (parent_id, path, "order", opening_day_id)
VALUES ($1, 'opening_days', $2, $3);
