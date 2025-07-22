SELECT id, COUNT(*) OVER () as total_count
FROM projects p
WHERE 
  p.not_eligible_at IS NULL
  AND ($1::uuid IS NULL OR p.round_id = $1)
  AND ($2::project_category_type IS NULL OR p.category = $2)
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;