BEGIN;

DROP TABLE IF EXISTS external_service_repos;
ALTER TABLE repo ADD COLUMN sources jsonb DEFAULT '{}'::jsonb NOT NULL;

COMMIT;
