ALTER TABLE users
ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();

ALTER TABLE organizations
ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();

ALTER TABLE organizations_members 
ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();

ALTER TABLE media
ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();

ALTER TABLE projects
ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();

ALTER TABLE oauth_connects 
ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();