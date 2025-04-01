ALTER TABLE kyb_verifications 
ADD CONSTRAINT kyb_verifications_user_org_unique 
UNIQUE (user_id, organization_id);