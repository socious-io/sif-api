SELECT c.*,
       COUNT(cc.id) AS children_count,
       COUNT(cl.id) AS likes,
       COALESCE(jsonb_agg(cr.reaction) FILTER (WHERE cr.reaction IS NOT NULL), '[]'::jsonb) AS reactions,
       EXISTS (SELECT 1 FROM comment_likes cl2 WHERE cl2.comment_id = c.id AND cl2.identity_id = $2) AS identity_liked,
       cr2.reaction AS identity_reaction,
       row_to_json(m.*) AS media,
       row_to_json(i.*) AS identity
FROM comments c
LEFT JOIN comments cc ON c.id = cc.parent_id
LEFT JOIN comment_likes cl ON c.id = cl.comment_id
LEFT JOIN comment_reactions cr ON c.id = cr.comment_id
LEFT JOIN comment_reactions cr2 ON c.id = cr2.comment_id AND identity_id = $2
LEFT JOIN media m ON c.media_id = m.id
JOIN identities i ON c.identity_id = i.id
WHERE c.parent_id=$1 AND c.identity_id=$2
LIMIT $3 OFFSET $4