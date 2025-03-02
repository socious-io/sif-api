CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TYPE status_type AS ENUM (
    'ACTIVE',
    'NOT_ACTIVE',
    'SUSPENDED'
);

CREATE TYPE project_status_type AS ENUM (
    'DRAFT',
    'EXPIRE',
    'ACTIVE'
);

CREATE TYPE identity_type AS ENUM (
    'users',
    'organizations'
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    city VARCHAR(255),
    country VARCHAR(255),
    address VARCHAR(255),
    avatar JSONB,
    cover JSONB,
    language VARCHAR(50),
    impact_points REAL DEFAULT 0,
    donates REAL DEFAULT 0,
    project_supported INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    shortname TEXT NOT NULL,
    name TEXT,
    bio TEXT,
    description TEXT,
    email TEXT,
    phone TEXT,
    
    city TEXT,
    country TEXT,
    address TEXT,
    website TEXT,
    
    mission TEXT,
    culture TEXT,
    
    logo JSONB,
    cover JSONB,
    
    status TEXT DEFAULT 'ACTIVE',
    
    verified_impact BOOLEAN NOT NULL DEFAULT FALSE,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at TIMESTAMP  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP  NOT NULL DEFAULT NOW()
);

CREATE TABLE identities (
    id UUID PRIMARY KEY,
    TYPE identity_type NOT NULL,
    meta JSONB,
    created_at TIMESTAMP  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP  NOT NULL DEFAULT NOW()
);


CREATE TABLE organizations_members(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);


CREATE TABLE media (
    id UUID PRIMARY KEY,
    identity_id UUID REFERENCES identities(id),
    url TEXT NOT NULL,
    filename TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE projects (
    id UUID PRIMARY KEY,
    title TEXT,
    description TEXT,
    status project_status_type NOT NULL DEFAULT 'ACTIVE',
    city TEXT,
    country TEXT,
    social_cause TEXT NOT NULL,
    identity_id UUID,
    cover_id UUID,
    wallet_address TEXT NOT NULL,
    wallet_env TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_identity FOREIGN KEY (identity_id) REFERENCES identities(id) ON DELETE CASCADE,
    CONSTRAINT fk_cover FOREIGN KEY (cover_id) REFERENCES media(id) ON DELETE SET NULL
);


CREATE TABLE oauth_connects (
    id UUID PRIMARY KEY,
    identity_id UUID REFERENCES identities(id),
    provider TEXT NOT NULL,
    matrix_unique_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    meta JSONB,
    expired_at TIMESTAMP,
    created_at TIMESTAMP  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP  NOT NULL DEFAULT NOW()
);

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
            'project_supported', NEW.project_supported
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

CREATE TRIGGER trigger_users_sync
AFTER INSERT OR UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION sync_identities();

CREATE TRIGGER trigger_organizations_sync
AFTER INSERT OR UPDATE ON organizations
FOR EACH ROW EXECUTE FUNCTION sync_identities();
