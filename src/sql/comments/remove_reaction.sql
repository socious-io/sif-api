DELETE FROM comment_reactions
WHERE comment_id = $1 AND identity_id = $2;
