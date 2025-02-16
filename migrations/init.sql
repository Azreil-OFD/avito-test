CREATE TABLE IF NOT EXISTS  users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT                NOT NULL,
    coins         INT     DEFAULT 0,
    is_deleted    BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS  inventory
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL,
    item_type  VARCHAR(255) NOT NULL,
    quantity   INT     DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS  transactions
(
    id          SERIAL PRIMARY KEY,
    sender_id   INT,
    receiver_id INT,
    amount      INT NOT NULL,
    timestamp   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted  BOOLEAN   DEFAULT FALSE,
    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS merch_store
(
    id         UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    item       VARCHAR(255) UNIQUE                 NOT NULL,
    price      DECIMAL(10, 2)                      NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);


INSERT INTO merch_store (item, price)
VALUES ('t-shirt', 80.00),
       ('cup', 20.00),
       ('book', 50.00),
       ('', 10.00),
       ('powerbank', 200.00),
       ('hoody', 300.00),
       ('umbrella', 200.00),
       ('socks', 10.00),
       ('wallet', 50.00),
       ('pink-hoody', 500.00)
ON CONFLICT (item) DO NOTHING;