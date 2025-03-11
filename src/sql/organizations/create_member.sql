INSERT INTO organizations_members (organization_id, user_id, created_at) VALUES
($1, $2, $3)
ON CONFLICT (organization_id, user_id) DO NOTHING