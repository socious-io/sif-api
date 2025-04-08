DELETE FROM comment_likes
WHERE comment_id = $1 AND identity_id = $2;
