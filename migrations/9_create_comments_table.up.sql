BEGIN;

CREATE TABLE IF NOT EXISTS comments (
    cid             TEXT                     NOT NULL,
    video_id        TEXT                     NOT NULL,
    text            TEXT                     NOT NULL,
    text_clean      TEXT                     NOT NULL,
    text_extra      jsonb,
    user_nickname   TEXT                     NOT NULL,
    user_avatar_uri TEXT                     NOT NULL,
    user_sec_uid    TEXT                     NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

ALTER TABLE ONLY comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY ( cid );

CREATE INDEX IF NOT EXISTS comments_video_id ON comments( video_id );

COMMIT;
