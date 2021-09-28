BEGIN;

ALTER TABLE videos
    DROP COLUMN IF EXISTS is_dynamic_cover;

DROP VIEW IF EXISTS video_list;

CREATE VIEW video_list AS (
                          SELECT v.*,
                                 members.sec_uid,
                                 members.uid,
                                 members.unique_id,
                                 members.short_id,
                                 members."name",
                                 members.avatar_url,
                                 members.signature
                          FROM (SELECT videos.*, array_agg(DISTINCT video_urls.url) AS video_urls
                                FROM videos JOIN video_urls ON video_urls.video_id = videos.id
                                GROUP BY videos.id) AS v JOIN members ON members.sec_uid = v.author_sec_id );
COMMIT;
