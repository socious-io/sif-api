INSERT INTO organizations (
    id, shortname, name, bio, description, email, phone,
    city, country, address, website, mission, culture,
    status, verified_impact, verified, logo, cover
) VALUES (
    $1, $2, $3, $4, $5, $6, 
    $7, $8, $9, $10, $11, $12, 
    $13, $14, $15, $16, $17, $18
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
    status = EXCLUDED.status,
    verified_impact = EXCLUDED.verified_impact,
    verified = EXCLUDED.verified,
    cover=EXCLUDED.cover,
    logo=EXCLUDED.logo
RETURNING id;