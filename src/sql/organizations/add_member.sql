INSERT INTO organizations_members (organization_id, user_id) VALUES ($1, $2)
ON CONFLICT (organization_id, user_id) DO NOTHING