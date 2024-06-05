CREATE TABLE IF NOT EXISTS user_favourite_episodes (
    id serial primary key,
    user_id serial,
    episode_id serial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    UNIQUE (user_id, episode_id),
    FOREIGN KEY (episode_id) REFERENCES podcast_episodes(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)