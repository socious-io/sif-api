SELECT r.*, 
       row_to_json(m.*) AS cover
FROM rounds r
LEFT JOIN media m ON m.id = r.cover_id
ORDER BY r.created_at DESC
LIMIT 1;

