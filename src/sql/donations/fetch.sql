SELECT d.*, row_to_json(u.*) AS user
FROM donations d 
JOIN users u ON d.user_id = u.id
JOIN unnest(ARRAY[?]::uuid[]) WITH ORDINALITY AS ids(id, ord) ON ids.id = d.id
ORDER BY ids.ord;