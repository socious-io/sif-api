CREATE UNIQUE INDEX unique_provider_mui
ON oauth_connects (matrix_unique_id, provider);