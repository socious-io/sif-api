SELECT p.*,
row_to_json(m.*) AS cover,
row_to_json(i.*) AS identity
FROM projects p
JOIN identities i ON i.id=p.identity_id
LEFT JOIN media m ON m.id=p.cover_id
WHERE p.id IN (?)
ORDER BY p.created_at DESC