INSERT INTO contracts (
  name,
  description,
  type,
  total_amount,
  currency,
  crypto_currency,
  currency_rate,
  commitment,
  commitment_period,
  commitment_period_count,
  payment_type,
  project_id,
  applicant_id,
  provider_id,
  client_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
ON CONFLICT (client_id, provider_id, project_id) WHERE status='CREATED' DO UPDATE SET id=contracts.id
RETURNING *