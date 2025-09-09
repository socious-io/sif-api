INSERT INTO projects (
    identity_id, title, description,
    status, cover_id, total_requested_amount,
    school_name, school_size, kpw,
    kwh_per_year, co2_per_year, expires_at
)
VALUES (
    $1, $2, $3,
    $4, $5, $6, 
    $7, $8, $9,
    $10, $11, $12
)
RETURNING *