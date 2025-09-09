-- donation paid_as
CREATE TYPE donation_paid_as_type AS ENUM('DONATION', 'INVESTMENT');
ALTER TABLE donations ADD COLUMN paid_as donation_paid_as_type DEFAULT 'DONATION';

-- project model change
ALTER TABLE projects
    --Removing
    DROP COLUMN search_vector,
    DROP COLUMN city,
    DROP COLUMN country,
    DROP COLUMN website,
    DROP COLUMN social_cause,
    DROP COLUMN round_id,
    DROP COLUMN wallet_address,
    DROP COLUMN wallet_env,
    DROP COLUMN linkedin,
    DROP COLUMN video,
    DROP COLUMN problem_statement,
    DROP COLUMN solution,
    DROP COLUMN goals,
    DROP COLUMN cost_beakdown,
    DROP COLUMN impact_assessment,
    DROP COLUMN impact_assessment_type,
    DROP COLUMN voluntery_contribution,
    DROP COLUMN feasibility,
    DROP COLUMN email,
    DROP COLUMN category,
    DROP COLUMN not_eligible_at,
    --Adding
    ADD COLUMN school_name text,
    ADD COLUMN school_size integer,
    ADD COLUMN kpw double precision,
    ADD COLUMN kwh_per_year double precision,
    ADD COLUMN co2_per_year double precision,
    ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B')
    ) STORED;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS investments REAL DEFAULT 0;

---------------

-- Drop triggers related to rounds/round_id
DROP TRIGGER IF EXISTS trigger_increment_total_projects ON projects;
DROP TRIGGER IF EXISTS trigger_decrement_total_projects ON projects;
DROP TRIGGER IF EXISTS trigger_increment_total_votes ON votes;
DROP TRIGGER IF EXISTS trigger_increment_total_donations ON donations;
DROP TRIGGER IF EXISTS trigger_set_latest_round ON projects;

-- Drop associated functions as well
DROP FUNCTION IF EXISTS increment_total_projects() CASCADE;
DROP FUNCTION IF EXISTS decrement_total_projects() CASCADE;
DROP FUNCTION IF EXISTS increment_total_votes() CASCADE;
DROP FUNCTION IF EXISTS increment_total_donations() CASCADE;
DROP FUNCTION IF EXISTS set_latest_round() CASCADE;


CREATE OR REPLACE FUNCTION update_user_total_investments()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET investments = investments + NEW.amount
    WHERE id = NEW.user_id AND paid_as='INVESTMENT';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_total_investments
AFTER INSERT ON donations
FOR EACH ROW EXECUTE FUNCTION update_user_total_investments();
