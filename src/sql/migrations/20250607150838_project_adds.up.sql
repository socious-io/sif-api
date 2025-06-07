CREATE TYPE project_category_type AS ENUM (
    'EMERGING_MARKETS',
    'OPEN_INNOVATION',
    'WOMEN_LEADERS'
);

ALTER TABLE projects 
    ADD COLUMN linkdin TEXT,
    ADD COLUMN video TEXT,
    ADD COLUMN problem_statement TEXT,
    ADD COLUMN solution TEXT,
    ADD COLUMN goals TEXT,
    ADD COLUMN total_requested_amount INTEGER DEFAULT 0,
    ADD COLUMN cost_beakdown TEXT,
    ADD COLUMN impact_assessment INTEGER DEFAULT 0,
    ADD COLUMN voluntery_contribution TEXT,
    ADD COLUMN feasibility TEXT,
    ADD COLUMN category project_category_type NOT NULL DEFAULT 'OPEN_INNOVATION';