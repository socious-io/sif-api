SELECT id, COUNT(*) OVER () as total_count
FROM donations
WHERE project_id=$1 AND status='APPROVED'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3