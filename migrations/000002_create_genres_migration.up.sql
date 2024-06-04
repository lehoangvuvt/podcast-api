CREATE TABLE IF NOT EXISTS genres (
    id serial primary key,
    uuid text NOT NULL UNIQUE,
    genre_name text NOT NULL,
    genre_desc text NOT NULL,
    bg_image text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone
)