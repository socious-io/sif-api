SELECT id, COUNT(*) OVER () as total_count
FROM projects p
ORDER BY p.created_at DESC
WHERE p.round_id = (SELECT id FROM rounds ORDER BY created_at DESC LIMIT 1)
LIMIT $1 OFFSET $2
