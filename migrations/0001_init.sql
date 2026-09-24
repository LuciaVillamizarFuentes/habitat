-- 0001_init.sql
-- Esquema inicial. Las tablas de tenants y units están listas; el resto lo diseñas tú (T-08, T-10, T-13, T-15).

CREATE TABLE IF NOT EXISTS units (
    id           TEXT PRIMARY KEY,
    code         TEXT        NOT NULL UNIQUE,
    kind         TEXT        NOT NULL CHECK (kind IN ('apartment','parking','storage','commerce')),
    floor        INT         NOT NULL,
    area_m2      NUMERIC(8,2) NOT NULL CHECK (area_m2 > 0),
    coefficient  NUMERIC(7,6) NOT NULL CHECK (coefficient > 0 AND coefficient <= 1),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tenants (
    id          TEXT PRIMARY KEY,
    full_name   TEXT        NOT NULL,
    email       TEXT        NOT NULL,
    phone       TEXT,
    unit_id     TEXT REFERENCES units(id),
    active      BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS tenants_email_lower_uq ON tenants (lower(email));

-- TODO(T-10): tabla resources y reservations.
--   Pista senior (T-17): mira EXCLUDE USING gist con tstzrange para que
--   Postgres rechace reservas solapadas incluso bajo concurrencia.
-- TODO(T-13): tabla tickets.
-- TODO(T-15): tabla audit_events (append-only).
