CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    debt_id bigint NOT NULL REFERENCES debts ON DELETE CASCADE,
    amount decimal(10, 2) NOT NULL,
    description text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version uuid NOT NULL DEFAULT uuid_generate_v4()
);