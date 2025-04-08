CREATE TABLE comments (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    project_id UUID NOT NULL,
    identity_id UUID NOT NULL,
    media_id UUID,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (identity_id) REFERENCES identities(id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE SET NULL
);


CREATE TABLE comment_likes (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    comment_id UUID NOT NULL,
    identity_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (identity_id) REFERENCES identities(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_comment_like ON comment_likes (identity_id, comment_id);

CREATE TABLE comment_reactions (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    comment_id UUID NOT NULL,
    identity_id UUID NOT NULL,
    reaction TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (identity_id) REFERENCES identities(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_comment_reaction ON comment_reactions (identity_id, comment_id);
