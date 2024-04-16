CREATE TABLE money_records (
                               id SERIAL PRIMARY KEY,
                               customer_id INTEGER NOT NULL REFERENCES customer ON DELETE CASCADE,
                               reference VARCHAR(50) UNIQUE NOT NULL,
                               status VARCHAR(50) NOT NULL,
                               amount DOUBLE PRECISION NOT NULL
);