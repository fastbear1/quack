-- +goose Up
CREATE INDEX IF NOT EXISTS "idx_ttb_float" ON "public"."test_table_b" USING hash (prec);
CREATE INDEX IF NOT EXISTS "idx_ttb_height_em" ON "public"."test_table_b" USING btree (height, em);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_test_table_b_deleted_at" ON "public"."test_table_b" USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS "idx_tta_date_birth" ON "public"."test_table_a" USING btree (date(birthday));
CREATE INDEX IF NOT EXISTS "idx_tta_where_sid" ON "public"."test_table_a" USING btree (sid);

-- +goose Down
DROP INDEX IF EXISTS "idx_ttb_float";
DROP INDEX IF EXISTS "idx_ttb_height_em";
DROP INDEX IF EXISTS "idx_test_table_b_deleted_at";
DROP INDEX IF EXISTS "idx_tta_date_birth";
DROP INDEX IF EXISTS "idx_tta_where_sid";
