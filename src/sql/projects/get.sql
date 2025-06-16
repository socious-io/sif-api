SELECT id, COUNT(*) OVER () as total_count
FROM projects p
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2
