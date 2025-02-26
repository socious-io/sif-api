SELECT 
    r.referred_by_id,
    u.wallet_address,
    (r.created_at > NOW() - INTERVAL '1 month') AS fee_discount
FROM referrings r
LEFT JOIN users u ON r.referred_by_id=u.id
WHERE
    r.referred_identity_id=$1 AND
    r.created_at > NOW() - INTERVAL '1 year' AND
    u.identity_verified