DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'impact_assessment_type') THEN
        CREATE TYPE impact_assessment_type AS ENUM ('OPTION_A', 'OPTION_B');
    END IF;
END$$;

ALTER TABLE projects
ADD COLUMN impact_assessment_type impact_assessment_type;