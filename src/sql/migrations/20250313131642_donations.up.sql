CREATE TYPE donation_status_type AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'RELEASED');

CREATE TABLE donations (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL,
    project_id UUID NOT NULL,
    currency TEXT NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    status donation_status_type NOT NULL DEFAULT 'PENDING',
    transaction_id TEXT,
    release_transaction_id TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

