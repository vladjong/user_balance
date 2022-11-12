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

CREATE TABLE services
(
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE orders
(
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE transactions
(
    id serial PRIMARY KEY,
    customer_id bigint REFERENCES customers (id) NOT NULL,
    service_id bigint REFERENCES services (id) NOT NULL,
    order_id bigint REFERENCES orders (id) NOT NULL,
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

CREATE TABLE expected_transactions
(
    id serial PRIMARY KEY,
    transaction_id bigint REFERENCES transactions(id) NOT NULL
);

CREATE VIEW history_report AS
SELECT h.id, s.name, t.cost, h.accounting_datetime
FROM history AS h
    JOIN transactions t ON t.id = h.transaction_id
    JOIN services s ON s.id = t.service_id;

CREATE VIEW customer_report AS
SELECT h.id, t.customer_id, s.name AS service_name, o.name AS order_name, t.cost AS sum, h.accounting_datetime as date
FROM history AS h
    JOIN transactions t ON t.id = h.transaction_id
    JOIN services s ON s.id = t.service_id
    JOIN orders o ON o.id = t.order_id;

INSERT INTO services
    VALUES (1, 'Упаковка'), (2, 'Доставка'), (3, 'Консультация'), (4, 'Пополнение');

INSERT INTO orders
    VALUES (1, 'А1'), (2, 'А2'), (3, 'А3'), (4, 'Баланс');
