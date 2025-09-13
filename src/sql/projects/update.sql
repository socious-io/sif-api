UPDATE projects p SET
    title=$2,
    description=$3,
    status=$4,
    cover_id=$5,
    total_requested_amount=$6,
    school_name=$7,
    school_size=$8,
    kpw=$9,
    kwh_per_year=$10,
    co2_per_year=$11,
    expires_at=$12
WHERE id=$1
RETURNING *