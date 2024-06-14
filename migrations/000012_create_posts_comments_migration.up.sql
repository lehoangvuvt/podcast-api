CREATE TABLE IF NOT EXISTS posts_comments (
    id serial primary key,
    post_id serial,
    user_id serial,
    reply_to_comment_id serial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    UNIQUE (post_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (reply_to_comment_id) REFERENCES posts_comments(id)
)