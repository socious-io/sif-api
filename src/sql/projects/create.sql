INSERT INTO projects (
    title,
    description,
    status,
    city,
    country,
    social_cause,
    identity_id,
    cover_id,
    wallet_address,
    wallet_env
)
VALUES (
    $1, $2, $3, 
    $4, $5, $6,
    $7, $8,
    $9, $10
)
RETURNING *