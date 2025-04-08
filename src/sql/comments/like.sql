INSERT INTO comment_likes (comment_id, identity_id)
VALUES ($1, $2)
ON CONFLICT (comment_id, identity_id) DO NOTHING
RETURNING *;
