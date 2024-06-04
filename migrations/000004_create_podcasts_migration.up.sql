CREATE TABLE IF NOT EXISTS podcasts (
    id serial primary key,
    uuid text NOT NULL UNIQUE,
    owner_id serial,
    podcast_name text NOT NULL,
    podcast_desc text NOT NULL,
    thumbnail_url text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    FOREIGN KEY (owner_id) REFERENCES users(id)
)