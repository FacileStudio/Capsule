-- Baseline: the schema as GORM AutoMigrate built it, captured from production
-- on 2026-08-06 with pg_dump --schema-only and stripped of environment noise
-- (SET preamble, \restrict markers, ownership, and the public. qualification so
-- the migration follows search_path rather than pinning one schema).
--
-- Production is stamped at this version rather than running it; see
-- Wiki/MIGRATIONS.md. It runs for real only against an empty database.

-- +goose Up
CREATE TABLE pastes (
    id character varying(24) NOT NULL,
    content text NOT NULL,
    burn_after_read boolean DEFAULT true NOT NULL,
    expires_at timestamp with time zone,
    max_views bigint,
    view_count bigint DEFAULT 0 NOT NULL,
    has_password boolean DEFAULT false NOT NULL,
    syntax character varying(50),
    delete_token character varying(64),
    burned boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone
);

ALTER TABLE ONLY pastes
    ADD CONSTRAINT pastes_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX idx_pastes_delete_token ON pastes USING btree (delete_token);

-- +goose Down
DROP TABLE pastes;
