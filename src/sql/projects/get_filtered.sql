SELECT 
    id, 
    COUNT(*) OVER () as total_count
FROM projects p
WHERE
    p.expires_at > NOW()
    AND ($1 = '' OR p.search_vector @@ to_tsquery('english', 
        (SELECT string_agg(word || ':*', ' & ') 
         FROM unnest(regexp_split_to_array($1, E'\\s+')) AS word
         WHERE word != '')
    ))
ORDER BY 
    CASE WHEN $1 != '' THEN ts_rank(search_vector, to_tsquery('english', 
        (SELECT string_agg(word || ':*', ' & ') 
         FROM unnest(regexp_split_to_array($1, E'\\s+')) AS word
         WHERE word != '')
    )) END DESC,
    p.created_at ASC
LIMIT $2 OFFSET $3;