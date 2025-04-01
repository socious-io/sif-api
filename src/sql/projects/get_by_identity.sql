SELECT id, COUNT(*) OVER () as total_count
FROM projects p
WHERE p.identity_id = $1 AND not_eligible_at IS NULL
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3