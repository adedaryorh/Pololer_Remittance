DROP TABLE IF EXISTS "transfer";
DROP TABLE IF EXISTS "entries";
ALTER TABLE "accounts" DROP CONSTRAINT "unique_customer_currency";
DROP TABLE IF EXISTS "accounts";
