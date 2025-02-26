SELECT c.*,
  row_to_json(id1.*) as provider,
  row_to_json(id2.*) as client,
  row_to_json(a.*) as applicant,
  row_to_json(p.*) as project,
  row_to_json(pay.*) as payment,
  EXISTS(SELECT f.id FROM feedbacks f WHERE f.contract_id=c.id AND f.identity_id=c.provider_id) AS provider_feedback,
  EXISTS(SELECT f.id FROM feedbacks f WHERE f.contract_id=c.id AND f.identity_id=c.client_id) AS client_feedback,
  (
    COALESCE(
      (SELECT
        jsonb_agg(
          json_build_object(
          'id', m.id,
          'url', m.url, 
          'filename', m.filename
          )
        )
      FROM media m
      LEFT JOIN contract_requirements_files crf ON crf.contract_id=c.id
      WHERE m.id = crf.document
      ),
      '[]'
    )
  ) AS requirement_files
FROM contracts c
JOIN identities id1 ON id1.id = c.provider_id
JOIN identities id2 ON id2.id = c.client_id
LEFT JOIN projects p ON p.id = c.project_id
LEFT JOIN applicants a ON a.id = c.applicant_id
LEFT JOIN gopay_payments pay ON pay.id = c.payment_id
WHERE c.id IN (?)
ORDER BY c.created_at DESC