INSERT INTO comments (project_id, identity_id, media_id, parent_id, content)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

