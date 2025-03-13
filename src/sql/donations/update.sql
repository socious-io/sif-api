UPDATE donations
SET status = $2, transaction_id = $3, release_transaction_id = $4
WHERE id = $1
RETURNING *;