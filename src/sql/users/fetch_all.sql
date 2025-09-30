SELECT
  id,
  COUNT(*) OVER() as total_count
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
