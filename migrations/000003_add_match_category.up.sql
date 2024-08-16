CREATE TABLE IF NOT EXISTS "match_category" (
    "id" SERIAL PRIMARY KEY,
    "l1" varchar(50) NOT NULL,
    "l2" varchar(50),
    "l3" varchar(50),
    "l4" varchar(50),
    "l5" varchar(50),
    "l6" varchar(50),
    "l7" varchar(50),
    "l8" varchar(50),
    "match_id" INTEGER
);

ALTER TABLE "match_category" ADD FOREIGN KEY ("match_id") REFERENCES "category" ("id");

ALTER TABLE "match_category" ADD CONSTRAINT unique_cat UNIQUE NULLS NOT DISTINCT  (l1, l2, l3, l4, l5, l6, l7, l8);