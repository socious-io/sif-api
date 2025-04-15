SELECT id, COUNT(*) OVER () as total_count
FROM comments c
WHERE project_id=$1
LIMIT $2 OFFSET $3