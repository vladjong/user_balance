CREATE TABLE customers
(
    id serial PRIMARY KEY,
    balance numeric NOT NULL
);

CREATE TABLE transactions
(
    id serial PRIMARY KEY,
    customer_id serial REFERENCES customers (id) NOT NULL,
    service_id serial NOT NULL,
    order_id serial NOT NULL,
    cost numeric NOT NULL,
    transaction_datetime timestamp NOT NULL
)

CREATE TABLE accounting
(
    id serial PRIMARY KEY,
    transaction_id serial REFERENCES transactions(id) NOT NULL,
    status_transaction boolean NOT NULL
)
