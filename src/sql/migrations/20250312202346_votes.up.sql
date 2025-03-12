ALTER TABLE users ADD COLUMN identity_verified_at TIMESTAMP;

CREATE TABLE votes (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL,
    project_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_vote ON votes (user_id, project_id);
