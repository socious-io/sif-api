INSERT INTO projects (
    identity_id, title, description, 
    payment_type, payment_scheme, payment_currency,
    payment_range_lower, payment_range_higher,
    experience_level, status,
    remote_preference, project_type, project_length,
    skills, causes_tags, country, city, geoname_id,
    other_party_id, other_party_title, other_party_url,
    expires_at, updated_at,
    weekly_hours_lower, weekly_hours_higher,
    commitment_hours_lower, commitment_hours_higher,
    job_category_id, kind, payment_mode
)
VALUES (
    $1, $2, $3, 
    $4, $5, $6,
    $7, $8,
    $9, (CASE WHEN $29::project_kind='SERVICE' THEN 'ACTIVE'::project_status ELSE $10::project_status END),
    $11, $12, $13,
    COALESCE($14, '{}'::text[]), COALESCE($15, '{}'::social_causes_type[]), $16, $17, $18,
    $19, $20, $21,
    $22, $23,
    $24, $25,
    $26, $27,
    $28, $29, $30
)
RETURNING *