UPDATE organizations
SET 
    status = COALESCE($2, status),
    verified_impact = $3,
    verified = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
