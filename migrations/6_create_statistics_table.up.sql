BEGIN;

CREATE TABLE IF NOT EXISTS statistics (
    id         TEXT                     NOT NULL,
    share      INT                      NOT NULL,
    forward    INT                      NOT NULL,
    digg       INT                      NOT NULL,
    play       INT                      NOT NULL,
    comment    INT                      NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS statistics_id_created_at ON statistics( id, created_at );

COMMIT;
