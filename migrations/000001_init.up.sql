CREATE TABLE customers
(
    id serial PRIMARY KEY,
    balance numeric(15, 2) NOT NULL
);

CREATE TABLE accounts
(
    id serial PRIMARY KEY,
    customer_id bigint REFERENCES customers (id) NOT NULL UNIQUE,
    balance numeric(15, 2) NOT NULL
);

CREATE TABLE transactions
(
    id serial PRIMARY KEY,
    customer_id bigint REFERENCES customers (id) NOT NULL,
    service_id bigint NOT NULL,
    order_id bigint NOT NULL,
    cost numeric(15, 2) NOT NULL,
    transaction_datetime timestamp NOT NULL
);

CREATE TABLE history
(
    id serial PRIMARY KEY,
    transaction_id bigint REFERENCES transactions(id) NOT NULL,
    accounting_datetime timestamp NOT NULL,
    status_transaction boolean NOT NULL
);
