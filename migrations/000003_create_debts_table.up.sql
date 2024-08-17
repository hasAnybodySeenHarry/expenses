CREATE TABLE IF NOT EXISTS debts (
    id bigserial PRIMARY KEY,
    lender_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    borrower_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    category citext UNIQUE NOT NULL,
    total_amount decimal(10, 2) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version uuid NOT NULL DEFAULT uuid_generate_v4()
);