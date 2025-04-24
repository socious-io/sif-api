UPDATE comments
SET content = $2, media_id = $3
WHERE id = $1
RETURNING *;
