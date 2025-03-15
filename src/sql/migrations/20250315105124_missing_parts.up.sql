ALTER TABLE rounds 
    ADD COLUMN total_projects INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN total_votes INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN total_donations NUMERIC(10, 2) NOT NULL DEFAULT 0,
    ADD COLUMN voting_announce_at TIMESTAMP;

CREATE OR REPLACE FUNCTION increment_total_projects()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE rounds
    SET total_projects = total_projects + 1
    WHERE id = NEW.round_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_increment_total_projects
AFTER INSERT ON projects
FOR EACH ROW EXECUTE FUNCTION increment_total_projects();

CREATE OR REPLACE FUNCTION decrement_total_projects()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE rounds
    SET total_projects = total_projects - 1
    WHERE id = OLD.round_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_decrement_total_projects
AFTER DELETE ON projects
FOR EACH ROW EXECUTE FUNCTION decrement_total_projects();

CREATE OR REPLACE FUNCTION increment_total_votes()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE rounds
    SET total_votes = total_votes + 1
    WHERE id = (
        SELECT round_id
        FROM projects
        WHERE id = NEW.project_id
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_increment_total_votes
AFTER INSERT ON votes
FOR EACH ROW EXECUTE FUNCTION increment_total_votes();

CREATE OR REPLACE FUNCTION increment_total_donations()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE rounds
    SET total_donations = total_donations + NEW.amount
    WHERE id = (
        SELECT round_id
        FROM projects
        WHERE id = NEW.project_id
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_increment_total_donations
AFTER INSERT ON donations
FOR EACH ROW EXECUTE FUNCTION increment_total_donations();

CREATE TYPE organization_status_type AS ENUM ('NOT_ACTIVE', 'PENDING', 'ACTIVE');

ALTER TABLE organizations DROP COLUMN status;
ALTER TABLE organizations ADD COLUMN status organization_status_type NOT NULL DEFAULT 'NOT_ACTIVE';

CREATE OR REPLACE FUNCTION update_user_total_donations()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET donates = donates + NEW.amount
    WHERE id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_total_donations
AFTER INSERT ON donations
FOR EACH ROW EXECUTE FUNCTION update_user_total_donations();

CREATE OR REPLACE FUNCTION sync_identities()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'users' THEN
        INSERT INTO identities (id, TYPE, meta, created_at, updated_at)
        VALUES (NEW.id, 'users', jsonb_build_object(
            'username', NEW.username,
            'first_name', NEW.first_name,
            'last_name', NEW.last_name,
            'email', NEW.email,
            'city', NEW.city,
            'country', NEW.country,
            'address', NEW.address,
            'avatar', NEW.avatar,
            'cover', NEW.cover,
            'language', NEW.language,
            'impact_points', NEW.impact_points,
            'donates', NEW.donates,
            'project_supported', NEW.project_supported,
            'identity_verified_at', NEW.identity_verified_at
        ), NOW(), NOW())
        ON CONFLICT (id) DO UPDATE
        SET meta = EXCLUDED.meta,
            updated_at = NOW();
    
    ELSIF TG_TABLE_NAME = 'organizations' THEN
        INSERT INTO identities (id, TYPE, meta, created_at, updated_at)
        VALUES (NEW.id, 'organizations', jsonb_build_object(
            'shortname', NEW.shortname,
            'name', NEW.name,
            'bio', NEW.bio,
            'description', NEW.description,
            'email', NEW.email,
            'phone', NEW.phone,
            'city', NEW.city,
            'country', NEW.country,
            'address', NEW.address,
            'website', NEW.website,
            'mission', NEW.mission,
            'culture', NEW.culture,
            'logo', NEW.logo,
            'cover', NEW.cover,
            'status', NEW.status,
            'verified_impact', NEW.verified_impact,
            'verified', NEW.verified
        ), NOW(), NOW())
        ON CONFLICT (id) DO UPDATE
        SET meta = EXCLUDED.meta,
            updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

UPDATE users SET id=id;
UPDATE organizations SET id=id;
