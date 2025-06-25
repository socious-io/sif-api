SELECT d.*, row_to_json(u.*) AS user
FROM donations d 
JOIN users u ON d.user_id = u.id
WHERE d.id IN (?)
