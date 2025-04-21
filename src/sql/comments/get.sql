SELECT c.*,
       (SELECT COUNT(*) FROM comments cc WHERE cc.parent_id = c.id) AS children_count,
       (SELECT COUNT(*) FROM comment_likes cl WHERE cl.comment_id = c.id) AS likes,
       (
           SELECT COALESCE(jsonb_agg(jsonb_build_object('reaction', cr.reaction, 'count', cr.count)), '[]'::jsonb)
           FROM (
               SELECT cr.reaction, COUNT(cr.reaction) AS count
               FROM comment_reactions cr
               WHERE cr.comment_id = c.id AND cr.reaction IS NOT NULL
               GROUP BY cr.reaction
           ) cr
       ) AS reactions,
       EXISTS (SELECT 1 FROM comment_likes cl2 WHERE cl2.comment_id = c.id AND cl2.identity_id = $2) AS identity_liked,
       COALESCE(cr2.reaction, '') AS identity_reaction,
       row_to_json(m.*) AS media,
       row_to_json(i.*) AS identity
FROM comments c
LEFT JOIN comment_reactions cr2 ON c.id = cr2.comment_id AND cr2.identity_id = $2
LEFT JOIN media m ON c.media_id = m.id
JOIN identities i ON c.identity_id = i.id
WHERE c.project_id=$1 AND c.parent_id IS NULL
LIMIT $3 OFFSET $4