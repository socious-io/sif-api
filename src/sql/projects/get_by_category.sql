SELECT id, COUNT(*) OVER () as total_count
FROM projects p
WHERE p.category = $1 AND p.not_eligible_at IS NULL
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3