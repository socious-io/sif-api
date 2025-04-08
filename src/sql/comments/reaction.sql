INSERT INTO comment_reactions (comment_id, identity_id, reaction)
VALUES ($1, $2, $3)
ON CONFLICT (comment_id, identity_id) DO UPDATE SET reaction = $3
RETURNING *;
