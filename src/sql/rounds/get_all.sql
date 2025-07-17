SELECT id, COUNT(*) OVER() AS total_count
FROM rounds
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;