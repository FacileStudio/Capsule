-- +goose Up
-- porte owns sign-in: identities, sessions and the one-time codes that carry a
-- browser login to the CLI. Pastes stay anonymous and reference none of this.
-- The tables are porte/pg Schema verbatim.

CREATE TABLE IF NOT EXISTS porte_users (
	id             bigserial PRIMARY KEY,
	facile_id      text UNIQUE,
	email          text NOT NULL UNIQUE,
	email_verified boolean NOT NULL DEFAULT false,
	name           text NOT NULL DEFAULT '',
	avatar_url     text NOT NULL DEFAULT '',
	avatar_source  text NOT NULL DEFAULT '',
	created_at     timestamptz NOT NULL DEFAULT now(),
	updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS porte_identities (
	user_id         bigint NOT NULL REFERENCES porte_users(id) ON DELETE CASCADE,
	provider        text NOT NULL,
	subject         text NOT NULL,
	password_hash   text NOT NULL DEFAULT '',
	access_token    text NOT NULL DEFAULT '',
	refresh_token   text NOT NULL DEFAULT '',
	token_expiry    timestamptz,
	roles           jsonb,
	roles_synced_at timestamptz,
	synced_at       timestamptz,
	created_at      timestamptz DEFAULT now(),
	PRIMARY KEY (provider, subject)
);
CREATE INDEX IF NOT EXISTS porte_identities_user_idx ON porte_identities (user_id);
ALTER TABLE porte_identities ADD COLUMN IF NOT EXISTS created_at timestamptz;
ALTER TABLE porte_identities ALTER COLUMN created_at SET DEFAULT now();

UPDATE porte_identities SET subject = user_id::text
 WHERE provider = 'local' AND subject <> user_id::text;

CREATE TABLE IF NOT EXISTS porte_sessions (
	id           bigserial PRIMARY KEY,
	token_hash   text NOT NULL UNIQUE,
	user_id      bigint NOT NULL REFERENCES porte_users(id) ON DELETE CASCADE,
	label        text NOT NULL DEFAULT '',
	created_at   timestamptz NOT NULL DEFAULT now(),
	last_used_at timestamptz NOT NULL DEFAULT now(),
	expires_at   timestamptz
);
CREATE INDEX IF NOT EXISTS porte_sessions_user_idx ON porte_sessions (user_id);
CREATE INDEX IF NOT EXISTS porte_sessions_expiry_idx ON porte_sessions (expires_at)
	WHERE expires_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS porte_login_codes (
	code_hash   text PRIMARY KEY,
	user_id     bigint NOT NULL REFERENCES porte_users(id) ON DELETE CASCADE,
	expires_at  timestamptz NOT NULL,
	consumed_at timestamptz
);
ALTER TABLE porte_login_codes ADD COLUMN IF NOT EXISTS consumed_at timestamptz;

-- +goose Down
DROP TABLE porte_login_codes;
DROP TABLE porte_sessions;
DROP TABLE porte_identities;
DROP TABLE porte_users;
