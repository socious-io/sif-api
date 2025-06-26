INSERT INTO donations (user_id, project_id, amount, currency, status, anonymous, rate)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;