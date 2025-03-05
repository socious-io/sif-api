INSERT INTO users (first_name, last_name, username, email, city, country, avatar, cover, language, impact_points) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;