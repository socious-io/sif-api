
DELETE FROM organizations_members
WHERE (user_id, organization_id) IN (
    SELECT user_id, organization_id
    FROM (
        SELECT 
            user_id, 
            organization_id,
            ROW_NUMBER() OVER (
                PARTITION BY user_id, organization_id 
                ORDER BY user_id, organization_id
            ) AS rn
        FROM organizations_members
    ) AS subquery
    WHERE rn > 1
);

CREATE UNIQUE INDEX unique_members ON organizations_members (user_id, organization_id);

