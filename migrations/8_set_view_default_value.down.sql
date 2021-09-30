BEGIN;

DROP VIEW IF EXISTS video_list;

CREATE VIEW video_list AS (
                          SELECT v.*,
                                 members.sec_uid,
                                 members.uid,
                                 members.unique_id,
                                 members.short_id,
                                 members."name",
                                 members.avatar_url,
                                 members.signature,
                                 statistics.share,
                                 statistics.forward,
                                 statistics.digg,
                                 statistics.play,
                                 statistics.comment
                          FROM (SELECT videos.*, array_agg(DISTINCT video_urls.url) AS video_urls
                                FROM videos JOIN video_urls ON video_urls.video_id = videos.id
                                GROUP BY videos.id) AS v
                               JOIN members ON members.sec_uid = v.author_sec_id
                               LEFT JOIN (SELECT *
                                          FROM statistics JOIN (SELECT id AS statistic_id, MAX(created_at) AS created_at
                                                                FROM statistics
                                                                GROUP BY id) AS latest_s
                                                               ON "statistics".id = latest_s.statistic_id AND
                                                                  "statistics".created_at =
                                                                  latest_s.created_at) AS statistics
                                         ON statistics.id = v.id );

COMMIT;
