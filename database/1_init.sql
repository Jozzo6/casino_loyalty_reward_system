\set ON_ERROR_STOP on

CREATE EXTENSION IF NOT EXISTS citext;

CREATE OR REPLACE FUNCTION update_modified_column()
	RETURNS trigger LANGUAGE plpgsql AS $function$
BEGIN
    NEW.updated = now();
    RETURN NEW; 
END;
$function$;

CREATE TABLE users (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	email CITEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	role INTEGER DEFAULT 0,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX users_email_idx ON users (email);

CREATE TRIGGER users_modtime BEFORE UPDATE
	ON users 
	FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE promotions (
	promotion_id UUID PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	amount NUMERIC NOT NULL,
	is_active BOOLEAN,
	start_date TIMESTAMPTZ NOT NULL,
	end_date TIMESTAMPTZ NOT NULL,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER promotions_modtime BEFORE UPDATE
	ON promotions 
	FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE users_promotions (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	promotion_id UUID REFERENCES promotions(promotion_id) ON DELETE CASCADE,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);