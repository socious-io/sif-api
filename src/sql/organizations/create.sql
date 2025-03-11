INSERT INTO organizations (
    id, shortname, name, bio, description, email, phone,
    city, country, address, website, mission, culture,
    logo, cover, status, verified_impact, verified
) VALUES (
    $1, $2, $3, $4, $5, $6, 
    $7, $8, $9, $10, $11, $12, 
    $13, $14, COALESCE($15, 'ACTIVE'), $16, $17, $18
)
ON CONFLICT (id)
DO UPDATE SET
    shortname = EXCLUDED.shortname,
    name = EXCLUDED.name,
    bio = EXCLUDED.bio,
    description = EXCLUDED.description,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    city = EXCLUDED.city,
    country = EXCLUDED.country,
    address = EXCLUDED.address,
    website = EXCLUDED.website,
    mission = EXCLUDED.mission,
    culture = EXCLUDED.culture,
    logo = EXCLUDED.logo,
    cover = EXCLUDED.cover,
    status = EXCLUDED.status,
    verified_impact = EXCLUDED.verified_impact,
    verified = EXCLUDED.verified
RETURNING id;