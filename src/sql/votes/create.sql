INSERT INTO votes (user_id, project_id)
VALUES ($1, $2)
RETURNING id;

