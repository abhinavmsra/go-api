CREATE TABLE IF NOT EXISTS "merchants"(
    id                  SERIAL PRIMARY KEY,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    name                varchar(255) NOT NULL UNIQUE,
    api_secret          varchar(255) NOT NULL UNIQUE
);

CREATE TRIGGER "set_public_merchants_updated_at"
BEFORE UPDATE ON "public"."merchants" 
FOR EACH ROW 
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();