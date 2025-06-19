CREATE TABLE accounts (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(100),
    user_email VARCHAR(100),
    password_hash VARCHAR(100) 
);