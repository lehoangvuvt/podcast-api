CREATE TABLE IF NOT EXISTS podcast_episodes (
    id serial primary key,
    uuid text NOT NULL UNIQUE,
    podcast_id serial,
    episode_name text NOT NULL,
    episode_no integer NOT NULL DEFAULT 1,
    episode_desc text NOT NULL,
    source_url text NOT NULL,
    created_at time with time zone NOT NULL DEFAULT NOW(),
    updated_at time with time zone,
    FOREIGN KEY (podcast_id) REFERENCES podcasts(id)
)