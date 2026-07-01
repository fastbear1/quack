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

-- +goose Down
DROP TABLE IF EXISTS "public"."simple_table";
DROP INDEX IF EXISTS "idx_simple_table_sid"
