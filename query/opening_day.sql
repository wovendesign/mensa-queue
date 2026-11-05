-- name: FindOrCreateOpeningDay :one
-- ON CONFLICT against opening_day_unique (NULLS NOT DISTINCT). The no-op
-- DO UPDATE makes the statement return the existing id on conflict.
INSERT INTO opening_day (day, opening_from, opening_to, service_from, service_to)
VALUES (@day::enum_opening_day_day, @opening_from, @opening_to, @service_from, @service_to)
ON CONFLICT (day, opening_from, opening_to, service_from, service_to)
    DO UPDATE SET day = EXCLUDED.day
RETURNING id;

-- name: FindOpeningDay :one
SELECT id
FROM opening_day
WHERE day = @day::enum_opening_day_day
  AND opening_from = @opening_from
  AND opening_to = @opening_to
  AND service_from IS NOT DISTINCT FROM @service_from
  AND service_to IS NOT DISTINCT FROM @service_to
LIMIT 1;
