SELECT id, COUNT(*) OVER () as total_count
FROM comments c
WHERE project_id=$1
LIMIT $3 OFFSET $4