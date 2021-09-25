BEGIN;

CREATE TABLE IF NOT EXISTS video_urls (
    video_id      TEXT                     NOT NULL,
    url           TEXT                     NOT NULL,
    status        TEXT                     NOT NULL,
    last_check_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

ALTER TABLE ONLY video_urls
    ADD CONSTRAINT video_urls_pkey PRIMARY KEY ( url );

CREATE INDEX video_urls_status ON video_urls( status );

CREATE INDEX video_urls_video_id ON video_urls( video_id, status );

COMMIT;
