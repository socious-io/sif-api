SELECT id, COUNT(*) OVER () as total_count
FROM comments c
WHERE project_id=$1 AND c.parent_id IS NULL
LIMIT $2 OFFSET $3