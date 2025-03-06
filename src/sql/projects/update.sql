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
    website=$11
WHERE id=$1
RETURNING *