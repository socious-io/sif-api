SELECT p.*,
    row_to_json(m.*) AS cover,
    row_to_json(i.*) AS identity,
    (
        SELECT jsonb_object_agg(currency, amount)
        FROM (
            SELECT d.currency, SUM(d.amount) AS amount
            FROM donations d
            WHERE d.project_id=p.id AND status='APPROVED' AND paid_as='DONATION'
            GROUP BY d.currency
        ) subquery
    ) AS total_donations,
    (
        SELECT jsonb_object_agg(currency, amount)
        FROM (
            SELECT d.currency, SUM(d.amount) AS amount
            FROM donations d
            WHERE d.project_id=p.id AND status='APPROVED' AND paid_as='INVESTMENT'
            GROUP BY d.currency
        ) subquery
    ) AS total_investments
FROM projects p
JOIN identities i ON i.id=p.identity_id
LEFT JOIN media m ON m.id=p.cover_id
WHERE p.id IN (?)
ORDER BY p.created_at ASC