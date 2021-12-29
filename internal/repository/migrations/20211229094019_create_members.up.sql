CREATE TABLE IF NOT EXISTS "members"(
    id                  SERIAL PRIMARY KEY,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    name                varchar(255) NOT NULL UNIQUE,
    api_secret          varchar(255) NOT NULL UNIQUE,
    merchant_id         integer NOT NULL, 
    FOREIGN KEY ("merchant_id") REFERENCES "public"."merchants"("id") ON UPDATE restrict ON DELETE restrict
);

CREATE TRIGGER "set_public_members_updated_at"
BEFORE UPDATE ON "public"."members" 
FOR EACH ROW 
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();