SELECT p.*,
row_to_json(jc.*) AS job_category,
row_to_json(i.*) AS identity,
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
		LEFT JOIN service_work_samples sws ON sws.service_id=p.id
		WHERE m.id = sws.document
		),
		'[]'
	)
) AS work_samples
FROM projects p
JOIN identities i ON i.id=p.identity_id
LEFT JOIN job_categories jc ON jc.id=p.job_category_id
WHERE p.id IN (?)
ORDER BY p.created_at DESC