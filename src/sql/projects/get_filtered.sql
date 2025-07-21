SELECT id, COUNT(*) OVER () as total_count
FROM projects p
WHERE 
  p.not_eligible_at IS NULL
  AND (COALESCE(NULLIF($1, '')::uuid, p.round_id) = p.round_id)
  AND (COALESCE(NULLIF($2, '')::project_category_type, p.category) = p.category)
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;