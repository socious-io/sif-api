CREATE TABLE rounds (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    
    name VARCHAR(255) NOT NULL,
    cover_id UUID REFERENCES media(id) ON DELETE SET NULL,
    pool_amount INT,
    
    voting_start_at TIMESTAMP NOT NULL CHECK (voting_start_at > submission_start_at),
    voting_end_at TIMESTAMP NOT NULL CHECK (voting_end_at > voting_start_at),
    
    submission_start_at TIMESTAMP NOT NULL,
    submission_end_at TIMESTAMP NOT NULL CHECK (submission_end_at > submission_start_at),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO rounds (name, pool_amount, voting_start_at, voting_end_at, submission_start_at, submission_end_at)
VALUES (
    'Round 1', 
    10000, 
    '2025-02-02 00:00:00', 
    '2025-03-01 00:00:00', 
    '2025-01-01 00:00:00', 
    '2025-02-01 00:00:00'    
);


ALTER TABLE projects
ADD COLUMN round_id UUID REFERENCES rounds(id) ON DELETE SET NULL DEFAULT (
    SELECT id FROM rounds ORDER BY created_at DESC LIMIT 1
);
