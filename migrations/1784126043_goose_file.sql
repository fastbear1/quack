-- +goose Up
CREATE TABLE "public"."simple_table"(
	id uuid NOT NULL default gen_random_uuid(),
	name varchar(255) NOT NULL,
	sid smallint NOT NULL,
	email varchar(255) NOT NULL,
	status varchar(10) NOT NULL default 'active'::text,
	name_t varchar(255) NOT NULL,
	created_at timestamp NOT NULL default now(),
	updated_at timestamp NOT NULL default now(),
	PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_simple_table_sid" ON "public"."simple_table" USING btree (sid);
CREATE TABLE "public"."clicks"(
	id uuid NOT NULL default gen_random_uuid(),
	created_at timestamp without time zone NOT NULL default now(),
	updated_at timestamp without time zone NOT NULL default now(),
	type text NOT NULL,
	user_id uuid NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "clicks_users_user_id_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);
DROP TABLE IF EXISTS "public"."auth_users";
DROP TABLE IF EXISTS "public"."recars";
ALTER TABLE "public"."cars" ADD COLUMN IF NOT EXISTS shifts smallint
ALTER TABLE "public"."commands" ADD COLUMN IF NOT EXISTS owner_id uuid NOT NULL
ALTER TABLE "public"."commands" DROP COLUMN IF EXISTS description
ALTER TABLE "public"."commands" ALTER COLUMN IF EXISTS cid SET DEFAULT uuidv4()
ALTER TABLE "public"."commands" ADD CONSTRAINT IF NOT EXISTS "commands_owner_owner_id_id" FOREIGN KEY ("owner_id") REFERENCES "public"."owners" ("id") ON DELETE CASCADE
DROP INDEX IF EXISTS "idx_commands_cid"
CREATE INDEX IF NOT EXISTS "idx_commands_cid" ON "public"."commands" USING btree (cid);

-- +goose Down
DROP TABLE IF EXISTS "public"."simple_table";
DROP INDEX IF EXISTS "idx_simple_table_sid"
DROP TABLE IF EXISTS "public"."clicks";
CREATE TABLE "public"."auth_users"(
	id uuid NOT NULL default gen_random_uuid(),
	user_id uuid NOT NULL,
	password text NOT NULL,
	created_at timestamp without time zone NOT NULL default now(),
	updated_at timestamp without time zone NOT NULL default now(),
	PRIMARY KEY ("id"),
	CONSTRAINT "fk_auth_users_users" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);
CREATE TABLE "public"."recars"(
	id uuid NOT NULL default gen_random_uuid(),
	name text NOT NULL,
	status text NOT NULL default 'active'::text,
	created_at timestamp without time zone NOT NULL default now(),
	updated_at timestamp without time zone NOT NULL default now(),
	PRIMARY KEY ("id")
);
ALTER TABLE "public"."cars" DROP COLUMN IF EXISTS shifts
ALTER TABLE "public"."commands" DROP COLUMN IF EXISTS owner_id
ALTER TABLE "public"."commands" ADD COLUMN IF NOT EXISTS description text
ALTER TABLE "public"."commands" DROP CONSTRAINT IF EXISTS "commands_owner_owner_id_id"
DROP INDEX IF EXISTS "idx_commands_cid"
CREATE INDEX IF NOT EXISTS "idx_commands_cid" ON "public"."" USING btree (cid);
