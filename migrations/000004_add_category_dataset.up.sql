CREATE TABLE IF NOT EXISTS "category_dataset" (
    "id" SERIAL PRIMARY KEY,
    "l1_in" varchar(50) NOT NULL,
    "l2_in" varchar(50),
    "l3_in" varchar(50),
    "l4_in" varchar(50),
    "l5_in" varchar(50),
    "l6_in" varchar(50),
    "l7_in" varchar(50),
    "l8_in" varchar(50),
    "full_path_out" varchar NOT NULL,
    "name_out" varchar NOT NULL,
    "version" varchar NOT NULL,
    "label" varchar NOT NULL
);
