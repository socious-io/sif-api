INSERT INTO job_categories(name, hourly_wage_dollars)
VALUES($1, $2)
RETURNING *;