SELECT p.*,
    row_to_json(m.*) AS cover,
    row_to_json(i.*) AS identity,
    row_to_json(r.*) AS round,
    (SELECT COUNT(*) FROM votes v WHERE v.project_id=p.id) AS total_votes,
    (
        SELECT jsonb_object_agg(currency, amount)
        FROM (
            SELECT d.currency, SUM(d.amount) AS amount, AVG(d.rate)
            FROM donations d
            WHERE d.project_id=p.id AND status='APPROVED'
            GROUP BY d.currency
        ) subquery
    ) AS total_donations
FROM projects p
JOIN identities i ON i.id=p.identity_id
LEFT JOIN media m ON m.id=p.cover_id
LEFT JOIN rounds r ON r.id=p.round_id
WHERE p.id IN (?) AND not_eligible_at IS NULL
ORDER BY p.created_at ASC