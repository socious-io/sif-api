INSERT INTO kyb_verifications(user_id, organization_id)
VALUES($1, $2)
RETURNING *