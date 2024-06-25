BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- DROP TABLE "db_migration" -----------------------------------
DROP TABLE IF EXISTS "public"."db_migration" CASCADE;
-- -------------------------------------------------------------

CREATE TABLE IF NOT EXISTS public.db_migration (
	"id" bigserial NOT NULL,
	"name" VARCHAR NOT NULL,
	"created_time" Timestamp Without Time Zone,
	PRIMARY KEY ( "id" ),
	CONSTRAINT "unique_db_migration_name" UNIQUE( "name" )
);

CREATE TABLE IF NOT EXISTS public.loan (
    id SERIAL PRIMARY KEY,
    borrower_id INT NOT NULL,
    principal_amount FLOAT NOT NULL,
    rate FLOAT NOT NULL,
    roi FLOAT NOT NULL,
    status int NOT NULL DEFAULT 1,
    agreement_letter_url TEXT,
	picture_proof_url TEXT,
	approver_id INT,
	approval_date TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS public.loan_investment (
	ID SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    investor_id VARCHAR NOT NULL,
    amount FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.loan_disbursement (
	ID SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    signed_agreement_url TEXT NOT NULL,
    disburser_employee_id VARCHAR NOT NULL,
    disbursement_date TIMESTAMPTZ NOT NULL
);
COMMIT;
