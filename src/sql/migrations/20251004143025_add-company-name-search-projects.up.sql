ALTER TABLE projects DROP COLUMN IF EXISTS search_vector;
DROP INDEX IF EXISTS idx_projects_search_vector;

ALTER TABLE projects ADD COLUMN IF NOT EXISTS search_vector tsvector;

CREATE OR REPLACE FUNCTION update_projects_search_vector() RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := (
        setweight(to_tsvector('english', COALESCE(LOWER(NEW.title), '')), 'A') ||
        setweight(to_tsvector('simple', COALESCE(LOWER((
            SELECT o.name
            FROM identities i
            JOIN organizations o ON i.id = o.id
            WHERE i.id = NEW.identity_id
        )), '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(LOWER(NEW.description), '')), 'C') ||
        setweight(to_tsvector('english', COALESCE(LOWER(NEW.social_cause), '')), 'D') ||
        setweight(to_tsvector('english', COALESCE(LOWER(NEW.city), '') || ' ' || COALESCE(LOWER(NEW.country), '')), 'D')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER projects_search_vector_trigger
BEFORE INSERT OR UPDATE
ON projects
FOR EACH ROW
EXECUTE FUNCTION update_projects_search_vector();

CREATE INDEX IF NOT EXISTS idx_projects_search_vector ON projects USING GIN(search_vector);

UPDATE projects SET title = title;