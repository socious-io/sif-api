INSERT INTO users (id, first_name, last_name, username, email, city, country, avatar, cover, language, impact_points, identity_verified_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (id) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    username = EXCLUDED.username,
    email = EXCLUDED.email,
    city = EXCLUDED.city,
    country = EXCLUDED.country,
    avatar = EXCLUDED.avatar,
    cover = EXCLUDED.cover,
    language = EXCLUDED.language,
    impact_points = EXCLUDED.impact_points,
    identity_verified_at = EXCLUDED.identity_verified_at
RETURNING *;
