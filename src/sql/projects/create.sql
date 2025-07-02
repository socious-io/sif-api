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
    wallet_env,
    website,
    linkedin,
    video,
    problem_statement,
    solution,
    goals,
    total_requested_amount,
    cost_beakdown,
    impact_assessment,
    impact_assessment_type,
    voluntery_contribution,
    feasibility,
    category,
    email
)
VALUES (
    $1, $2, $3, 
    $4, $5, $6,
    $7, $8, $9,
    $10, $11, $12,
    $13, $14, $15,
    $16, $17, $18,
    $19, $20, $21,
    $22, $23, $24
)
RETURNING *