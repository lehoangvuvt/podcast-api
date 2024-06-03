CREATE TABLE IF NOT EXISTS genres_podcasts (
    id serial primary key,
    genre_id serial,
    podcast_id serial,
    created_at time with time zone NOT NULL DEFAULT NOW(),
    updated_at time with time zone,
    UNIQUE (genre_id, podcast_id),
    FOREIGN KEY (genre_id) REFERENCES genres(id),
    FOREIGN KEY (podcast_id) REFERENCES podcasts(id)
)