CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "username" varchar(30) UNIQUE NOT NULL,
    "password" varchar(60) NOT NULL,
    "active" bool NOT NULL DEFAULT true
);
