CREATE TABLE IF NOT EXISTS posts (
    id serial primary key,
    user_id serial,
    slug text NOT NULL UNIQUE,
    title text NOT NULL,
    short_content text NOT NULL,
    thumbnail_url text,
    content text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    FOREIGN KEY (user_id) REFERENCES users(id)
)