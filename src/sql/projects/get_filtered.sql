SELECT 
    id, 
    COUNT(*) OVER () as total_count
FROM projects p
WHERE 
    p.not_eligible_at IS NULL
    AND (COALESCE(NULLIF($1, '')::uuid, p.round_id) = p.round_id)
    AND (COALESCE(NULLIF($2, '')::project_category_type, p.category) = p.category)
    AND ($3 = '' OR p.search_vector @@ to_tsquery('english', 
        (SELECT string_agg(word || ':*', ' & ') 
         FROM unnest(regexp_split_to_array($3, E'\\s+')) AS word
         WHERE word != '')
    ))
ORDER BY 
    CASE WHEN $3 != '' THEN ts_rank(search_vector, to_tsquery('english', 
        (SELECT string_agg(word || ':*', ' & ') 
         FROM unnest(regexp_split_to_array($3, E'\\s+')) AS word
         WHERE word != '')
    )) END DESC,
    p.created_at ASC
LIMIT $4 OFFSET $5;