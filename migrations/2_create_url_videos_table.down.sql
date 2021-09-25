BEGIN;

DROP TABLE IF EXISTS video_urls;

ALTER TABLE videos
    ADD COLUMN video_urls TEXT[] NOT NULL;

ALTER TABLE videos
    ADD COLUMN video_cdn_url TEXT NOT NULL;

COMMIT;
