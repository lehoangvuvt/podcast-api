CREATE TABLE IF NOT EXISTS topics (
    id serial primary key,
    slug text NOT NULL UNIQUE,
    topic_name text NOT NULL UNIQUE,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone
)