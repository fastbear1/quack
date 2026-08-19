-- +goose Up
ALTER TABLE "public"."test_table_b" ADD COLUMN IF NOT EXISTS tt_a uuid NOT NULL;
ALTER TABLE "public"."test_table_b" ADD CONSTRAINT "fk_TestTableB_tt_a__test_table_a_id" FOREIGN KEY (tt_a) REFERENCES "public"."test_table_a" (id) ON DELETE CASCADE ON UPDATE SET NULL;
ALTER TABLE "public"."test_table_c" ADD COLUMN IF NOT EXISTS tt_b bigint NOT NULL;
ALTER TABLE "public"."test_table_c" ADD CONSTRAINT "ref_ttc_ttb_id" FOREIGN KEY (tt_b) REFERENCES "public"."test_table_b" (id) ON UPDATE SET NULL;

-- +goose Down
ALTER TABLE "public"."test_table_b" DROP COLUMN IF EXISTS tt_a;
ALTER TABLE "public"."test_table_b" DROP CONSTRAINT IF EXISTS "fk_TestTableB_tt_a__test_table_a_id;"
ALTER TABLE "public"."test_table_c" DROP COLUMN IF EXISTS tt_b;
ALTER TABLE "public"."test_table_c" DROP CONSTRAINT IF EXISTS "ref_ttc_ttb_id;"
