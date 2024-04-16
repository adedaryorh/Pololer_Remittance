CREATE TABLE "accounts" (
                           id BIGSERIAL PRIMARY KEY,
                           customer_id INTEGER NOT NULL,
                           balance DOUBLE PRECISION NOT NULL DEFAULT 0,
                           account_type varchar(256) NOT NULL,
                           account_status varchar(256) NOT NULL,
                           currency VARCHAR(10) NOT NULL,
                           created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           FOREIGN KEY (customer_id) REFERENCES customer(id) ON DELETE CASCADE
);

ALTER TABLE "accounts" ADD CONSTRAINT "unique_customer_currency"
UNIQUE (customer_id, currency);


CREATE TABLE "entries"(
                          id BIGSERIAL PRIMARY KEY,
                          account_id INTEGER NOT NULL,
                          amount DOUBLE PRECISION NOT NULL,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE "transfers"(
                            id BIGSERIAL PRIMARY KEY,
                            from_account_id INTEGER NOT NULL,
                            to_account_id INTEGER NOT NULL,
                            amount DOUBLE PRECISION NOT NULL,
                            created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (from_account_id) REFERENCES customer(id) ON DELETE CASCADE,
                            FOREIGN KEY (to_account_id) REFERENCES customer(id) ON DELETE CASCADE
);
