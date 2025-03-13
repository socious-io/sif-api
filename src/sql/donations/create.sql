INSERT INTO donations (user_id, project_id, amount, currency, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;