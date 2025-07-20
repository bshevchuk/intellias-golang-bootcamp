CREATE TABLE "feeds"(
    "id" integer GENERATED ALWAYS AS IDENTITY,
    "url" text,
    "createdAt" timestamp DEFAULT 'now()',
    PRIMARY KEY ("id")
)

