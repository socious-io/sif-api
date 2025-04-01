SELECT *
FROM votes
WHERE user_id = $1 AND project_id = $2
ORDER BY created_at DESC;