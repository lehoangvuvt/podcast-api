CREATE TABLE IF NOT EXISTS posts_comments (
    id serial primary key,
    post_id serial,
    user_id serial,
    reply_to_comment_id int DEFAULT (NULL),
    content text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (reply_to_comment_id) REFERENCES posts_comments(id)
)