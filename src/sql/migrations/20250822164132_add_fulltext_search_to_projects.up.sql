ALTER TABLE projects 
ADD COLUMN IF NOT EXISTS search_vector tsvector 
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', COALESCE(title, '')), 'A') || 
    setweight(to_tsvector('english', COALESCE(description, '')), 'B') || 
    setweight(to_tsvector('english', COALESCE(social_cause, '')), 'C') || 
    setweight(to_tsvector('english', COALESCE(city, '') || ' ' || COALESCE(country, '')), 'D')
) STORED;

CREATE INDEX IF NOT EXISTS idx_projects_search_vector ON projects USING GIN(search_vector);