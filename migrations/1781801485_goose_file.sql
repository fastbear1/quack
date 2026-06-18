-- +goose Up
CREATE TABLE "public"."simple_table"(
	id uuid PRIMARY KEY NOT NULL default gen_random_uuid(),
	name varchar(255) NOT NULL,
	sid smallint NOT NULL,
	email varchar(255) NOT NULL,
	status varchar(10) NOT NULL default active,
	name_t varchar(255) NOT NULL,
	created_at timestamp NOT NULL default now(),
	updated_at timestamp NOT NULL default now()
);
-- +goose Down
DROP TABLE IF EXISTS "public"."simple_table";