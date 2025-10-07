ALTER TABLE donations
    ALTER COLUMN rate DROP DEFAULT,
    ALTER COLUMN rate TYPE double precision
    USING rate::double precision,
    ALTER COLUMN rate SET DEFAULT 1;

UPDATE donations SET rate=0.0066 WHERE currency='JPY';
UPDATE donations SET rate=0.85 WHERE currency='lovelace' OR currency='ADA';