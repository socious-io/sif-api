SELECT r.*,
       row_to_json(m.*) AS cover
FROM rounds r
LEFT JOIN media m ON m.id = r.cover_id
WHERE r.id IN (?)
ORDER BY r.created_at DESC