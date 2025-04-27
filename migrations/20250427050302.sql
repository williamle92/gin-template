-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "email" SET NOT NULL, ALTER COLUMN "phone_number" SET NOT NULL, ADD CONSTRAINT "uni_users_phone_number" UNIQUE ("phone_number");
-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_phone_number" to table: "users"
CREATE INDEX "idx_users_phone_number" ON "public"."users" ("phone_number");
