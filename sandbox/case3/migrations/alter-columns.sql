-- +goose Up
ALTER TABLE "public"."test_table_a" ALTER COLUMN sid SET DEFAULT 100;
ALTER TABLE "public"."test_table_a" ALTER COLUMN age TYPE bigint;
ALTER TABLE "public"."test_table_a" ALTER COLUMN tags SET DEFAULT '{}'::text[];
ALTER TABLE "public"."test_table_a" ALTER COLUMN context TYPE varchar(255);
ALTER TABLE "public"."test_table_a" ALTER COLUMN context SET NOT NULL;
ALTER TABLE "public"."test_table_b" ALTER COLUMN em TYPE smallint;
ALTER TABLE "public"."test_table_b" ALTER COLUMN em DROP NOT NULL;
ALTER TABLE "public"."test_table_b" ALTER COLUMN prec SET NOT NULL;
ALTER TABLE "public"."test_table_b" ALTER COLUMN prec SET DEFAULT  0.002;

-- +goose Down
ALTER TABLE "public"."test_table_a" ALTER COLUMN sid SET DEFAULT 0;
ALTER TABLE "public"."test_table_a" ALTER COLUMN age TYPE smallint;
ALTER TABLE "public"."test_table_a" ALTER COLUMN tags DROP DEFAULT;
ALTER TABLE "public"."test_table_a" ALTER COLUMN context TYPE varchar(100);
ALTER TABLE "public"."test_table_a" ALTER COLUMN context DROP NOT NULL;
ALTER TABLE "public"."test_table_b" ALTER COLUMN em TYPE bigint;
ALTER TABLE "public"."test_table_b" ALTER COLUMN em SET NOT NULL;
ALTER TABLE "public"."test_table_b" ALTER COLUMN prec DROP NOT NULL;
ALTER TABLE "public"."test_table_b" ALTER COLUMN prec DROP DEFAULT;
