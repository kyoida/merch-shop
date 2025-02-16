CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    coins INT NOT NULL DEFAULT 1000,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user_id INT,
    to_user_id INT NOT NULL,
    amount INT NOT NULL,
    type TEXT NOT NULL,  -- 'transfer' or 'purchase'
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_id) REFERENCES users(id)
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    item TEXT NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
