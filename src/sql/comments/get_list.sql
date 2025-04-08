SELECT id, COUNT(*) OVER () as total_count
FROM comments c
WHERE project_id=$1 AND identity_id=$2
LIMIT $3 OFFSET $4