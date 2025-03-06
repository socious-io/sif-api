UPDATE organizations
SET 
    shortname = $2,
    name = $3,
    bio = $4,
    description = $5,
    email = $6,
    phone = $7,
    city = $8,
    country = $9,
    address = $10,
    website = $11,
    mission = $12,
    culture = $13,
    logo = $14,
    cover = $15,
    status = COALESCE($16, status),
    verified_impact = $17,
    verified = $18,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
