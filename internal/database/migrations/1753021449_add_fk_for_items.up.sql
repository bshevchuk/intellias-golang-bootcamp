ALTER TABLE "items"
    ADD COLUMN "feedId" integer,
    ADD FOREIGN KEY ("feedId") REFERENCES "feeds"("id") ON DELETE CASCADE;
