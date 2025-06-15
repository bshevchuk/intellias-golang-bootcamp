CREATE TABLE "public"."items" (
    "id" integer GENERATED ALWAYS AS IDENTITY,
    "title" text,
    "link" text,
    "description" text,
    "createdAt" timestamp DEFAULT 'now()',
    PRIMARY KEY ("id")
);
