CREATE TYPE user_role_type AS ENUM (
    'PUBLIC',
    'ADMIN', 
    'SCHOOL',
    'CEB',
    'INSTALLER',
    'GOVERNMENT'
);

ALTER TABLE users 
ADD COLUMN role user_role_type NOT NULL DEFAULT 'PUBLIC';
