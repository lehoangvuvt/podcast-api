CREATE TABLE IF NOT EXISTS posts_likes (
    id serial primary key,
    post_id serial,
    user_id serial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    UNIQUE (post_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
)