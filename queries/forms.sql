-- forms
-- name: get-forms
SELECT * FROM forms ORDER BY created_at;

-- name: get-form
SELECT * FROM forms WHERE (CASE WHEN $1 > 0 THEN id = $1 WHEN $2 != '' THEN uuid = $2::UUID END);

-- name: create-form
INSERT INTO forms (uuid, name, list_ids, redirect_url) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: update-form
UPDATE forms SET
    name = (CASE WHEN $2 != '' THEN $2 ELSE name END),
    list_ids = $3::INT[],
    redirect_url = $4,
    updated_at = NOW()
WHERE id = $1 RETURNING id;

-- name: delete-form
DELETE FROM forms WHERE id = $1;
