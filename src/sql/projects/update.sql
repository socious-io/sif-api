UPDATE projects p SET
    title=$2,
    description=$3,
    status=$4,
    city=$5,
    country=$6,
    social_cause=$7,
    cover_id=$8,
    wallet_address=$9,
    wallet_env=$10,
    website=$11,
    linkdin=$12,
    video=$13,
    problem_statement=$14,
    solution=$15,
    goals=$16,
    total_requested_amount=$17,
    cost_beakdown=$18,
    impact_assessment=$19,
    voluntery_contribution=$20,
    feasibility=$21,
    category=$22
WHERE id=$1
RETURNING *