ALTER TABLE "movements" DROP COLUMN "transfer_id";
ALTER TABLE "transfers" DROP CONSTRAINT "transfers_from_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT "transfers_to_account_id_fkey";

DROP TABLE IF EXISTS "transfers";