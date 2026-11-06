-- +goose Up
-- porte v0.10 binds a login code to the identity that vouched for it, stores the
-- person's locale, and v0.12 erases magic links with the rest of an identity.
-- Capsule never mounts magic links, but EraseIdentity deletes from the table,
-- so it has to exist. deletion_state holds the GDPR receiver's cursor.

ALTER TABLE porte_users ADD COLUMN IF NOT EXISTS locale text NOT NULL DEFAULT 'en';
ALTER TABLE porte_login_codes
	ADD COLUMN IF NOT EXISTS provider text NOT NULL DEFAULT '',
	ADD COLUMN IF NOT EXISTS subject text NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS porte_magic_links (
	token_hash  text PRIMARY KEY,
	user_id     bigint NOT NULL REFERENCES porte_users(id) ON DELETE CASCADE,
	ip          text NOT NULL DEFAULT '',
	created_at  timestamptz NOT NULL,
	expires_at  timestamptz NOT NULL,
	consumed_at timestamptz
);
CREATE INDEX IF NOT EXISTS porte_magic_links_user_idx ON porte_magic_links (user_id, created_at);
CREATE INDEX IF NOT EXISTS porte_magic_links_ip_idx ON porte_magic_links (ip, created_at)
	WHERE ip <> '';
CREATE INDEX IF NOT EXISTS porte_magic_links_expiry_idx ON porte_magic_links (expires_at);

CREATE TABLE IF NOT EXISTS deletion_state (
	id         smallint PRIMARY KEY CHECK (id = 1),
	state      bytea NOT NULL,
	updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE deletion_state;
DROP TABLE porte_magic_links;
ALTER TABLE porte_login_codes DROP COLUMN IF EXISTS subject, DROP COLUMN IF EXISTS provider;
ALTER TABLE porte_users DROP COLUMN IF EXISTS locale;
