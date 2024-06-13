CREATE TABLE IF NOT EXISTS posts_topics (
    id serial primary key,
    post_id serial,
    topic_id serial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    UNIQUE (post_id, topic_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (topic_id) REFERENCES topics(id)
)