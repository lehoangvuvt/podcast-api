CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    username text NOT NULL UNIQUE,
    hashed_password text NOT NULL,
    email text NOT NULL UNIQUE,
    is_active boolean DEFAULT false,
    created_at time with time zone NOT NULL DEFAULT NOW(),
    updated_at time with time zone
)