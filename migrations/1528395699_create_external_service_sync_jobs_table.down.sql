BEGIN;

DROP TABLE IF EXISTS external_service_sync_jobs;
DROP SEQUENCE IF EXISTS external_service_sync_jobs_id_seq;
DROP VIEW IF EXISTS external_service_sync_jobs_with_next_sync_at;

COMMIT;
