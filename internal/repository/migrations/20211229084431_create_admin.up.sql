CREATE TABLE IF NOT EXISTS "admins"(
    id                  SERIAL PRIMARY KEY,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    name                varchar(255) NOT NULL UNIQUE,
    api_secret          varchar(255) NOT NULL UNIQUE
);

CREATE TRIGGER "set_public_admins_updated_at"
BEFORE UPDATE ON "public"."admins" 
FOR EACH ROW 
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();