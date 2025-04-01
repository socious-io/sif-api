SELECT id, COUNT(*) OVER () as total_count
FROM projects p
WHERE p.round_id = (SELECT id FROM rounds ORDER BY created_at DESC LIMIT 1) AND not_eligible_at IS NULL
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2
