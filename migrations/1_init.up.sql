BEGIN;

CREATE TABLE IF NOT EXISTS members (
    sec_uid    TEXT                     NOT NULL,
    uid        TEXT                     NOT NULL,
    unique_id  TEXT                     NOT NULL,
    short_id   TEXT                     NOT NULL,
    name       TEXT                     NOT NULL,
    avatar_url TEXT                     NOT NULL,
    signature  TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_member_sec_uid ON members( sec_uid );

CREATE TABLE IF NOT EXISTS videos (
    id                 TEXT                     NOT NULL,
    author_sec_id      TEXT                     NOT NULL,
    description        TEXT                     NOT NULL,
    text_extra         TEXT[]                   NOT NULL,
    origin_cover_urls  TEXT[]                   NOT NULL,
    dynamic_cover_urls TEXT[]                   NOT NULL,
    video_height       INT                      NOT NULL,
    video_width        INT                      NOT NULL,
    video_duration     BIGINT                   NOT NULL,
    video_ratio        TEXT                     NOT NULL,
    video_urls         TEXT[]                   NOT NULL,
    video_cdn_url      TEXT                     NOT NULL,
    created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

ALTER TABLE ONLY videos
    ADD CONSTRAINT videos_pkey PRIMARY KEY ( id );

COMMIT;
