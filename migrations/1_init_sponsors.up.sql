BEGIN;

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS admins (
    id serial PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    acl text ARRAY,
    created_at timestamp DEFAULT NOW() NOT NULL,
    updated_at timestamp DEFAULT NOW() NOT NULL
);

COMMIT;