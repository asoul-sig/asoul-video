BEGIN;

ALTER TABLE videos
    DROP COLUMN IF EXISTS cover_width;

ALTER TABLE videos
    DROP COLUMN IF EXISTS cover_height;

COMMIT;
