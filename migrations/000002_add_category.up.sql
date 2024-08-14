CREATE TABLE IF NOT EXISTS "category" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(50) NOT NULL,
    "parent_id" INTEGER
);

ALTER TABLE "category" ADD FOREIGN KEY ("parent_id") REFERENCES "category" ("id");