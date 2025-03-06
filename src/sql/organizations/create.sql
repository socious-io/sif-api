INSERT INTO organizations (
    shortname, name, bio, description, email, phone,
    city, country, address, website, mission, culture,
    logo, cover, status, verified_impact, verified
) VALUES (
    $1, $2, $3, $4, $5, $6, 
    $7, $8, $9, $10, $11, $12, 
    $13, $14, COALESCE($15, 'ACTIVE'), $16, $17
) RETURNING id;