-- +goose Up
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS email varchar(255);
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS age smallint NOT NULL;
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS context varchar(100);
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS sid bigint NOT NULL DEFAULT 0;
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS salt text;
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS tags text[] NOT NULL;
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS name text NOT NULL;
ALTER TABLE "public"."test_table_a" ADD COLUMN IF NOT EXISTS birthday timestamp without time zone;
ALTER TABLE "public"."test_table_b" ADD COLUMN IF NOT EXISTS prec real;
ALTER TABLE "public"."test_table_b" ADD COLUMN IF NOT EXISTS height smallint;
ALTER TABLE "public"."test_table_b" ADD COLUMN IF NOT EXISTS em bigint NOT NULL;
ALTER TABLE "public"."test_table_b" ADD COLUMN IF NOT EXISTS sid uuid NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE "public"."test_table_c" ADD COLUMN IF NOT EXISTS data json NOT NULL DEFAULT '{}'::json;
ALTER TABLE "public"."test_table_c" ADD COLUMN IF NOT EXISTS tasks jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE "public"."test_table_c" ADD COLUMN IF NOT EXISTS colors integer[] NOT NULL DEFAULT '{}'::integer[];

-- +goose Down
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS email;
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS age;
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS context;
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS sid;
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS salt;
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS tags;
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS name;
ALTER TABLE "public"."test_table_a" DROP COLUMN IF EXISTS birthday;
ALTER TABLE "public"."test_table_b" DROP COLUMN IF EXISTS prec;
ALTER TABLE "public"."test_table_b" DROP COLUMN IF EXISTS height;
ALTER TABLE "public"."test_table_b" DROP COLUMN IF EXISTS em;
ALTER TABLE "public"."test_table_b" DROP COLUMN IF EXISTS sid;
ALTER TABLE "public"."test_table_c" DROP COLUMN IF EXISTS data;
ALTER TABLE "public"."test_table_c" DROP COLUMN IF EXISTS tasks;
ALTER TABLE "public"."test_table_c" DROP COLUMN IF EXISTS colors;
