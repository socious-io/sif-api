UPDATE users
SET role = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
