ALTER TABLE IF EXISTS "match_category" DROP CONSTRAINT IF EXISTS "match_category_match_id_fkey";

ALTER TABLE IF EXISTS "match_category" DROP CONSTRAINT IF EXISTS "unique_cat";

DROP TABLE IF EXISTS "match_category";