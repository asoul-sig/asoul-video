BEGIN;

ALTER TABLE videos
    ADD COLUMN face_points jsonb DEFAULT '[]'::jsonb NOT NULL;

COMMIT;
