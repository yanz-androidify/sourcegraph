BEGIN;

DROP TRIGGER IF EXISTS trig_soft_delete_repo_reference_on_external_service_repos ON repo;
DROP TRIGGER IF EXISTS trig_soft_delete_external_service_reference_on_external_service_repos ON repo;

DROP FUNCTION IF EXISTS soft_delete_external_service_reference_on_external_service_repos;
DROP FUNCTION IF EXISTS soft_delete_repo_reference_on_external_service_repos;

DROP TABLE IF EXISTS external_service_repos;
ALTER TABLE repo ADD COLUMN sources jsonb DEFAULT '{}'::jsonb NOT NULL;

COMMIT;
