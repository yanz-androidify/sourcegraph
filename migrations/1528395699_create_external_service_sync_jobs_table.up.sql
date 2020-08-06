BEGIN;

CREATE SEQUENCE external_service_sync_jobs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS external_service_sync_jobs (
    -- Columns required by workerutil.Store
    id integer NOT NULL DEFAULT nextval('external_service_sync_jobs_id_seq'::regclass),
    state text NOT NULL DEFAULT 'queued'::text,
    failure_message text,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    process_after timestamp with time zone,
    num_resets integer not null DEFAULT 0,
    -- Extra columns
    external_service_id bigint,
    -- Constraints
    CONSTRAINT external_services_id_fk
    FOREIGN KEY(id)
    REFERENCES external_services(id)
);

-- NOTE: No index on the state column was added as we expect the size of this table to stay fairly small

COMMIT;
