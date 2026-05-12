CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
   id BIGSERIAL PRIMARY KEY,
   username VARCHAR(50) UNIQUE NOT NULL,
   email VARCHAR(255) UNIQUE NOT NULL,
   password_hash TEXT NOT NULL,
   created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE accounts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  balance NUMERIC(15,2) DEFAULT 0,
  currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE cards (
   id BIGSERIAL PRIMARY KEY,
   account_id BIGINT REFERENCES accounts(id),
   encrypted_pan BYTEA NOT NULL,
   encrypted_expiry BYTEA NOT NULL,
   cvv_hash TEXT NOT NULL,
   pan_hmac TEXT NOT NULL,
   created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE transactions (
      id BIGSERIAL PRIMARY KEY,
      from_account_id BIGINT,
      to_account_id BIGINT,
      amount NUMERIC(15,2),
      type VARCHAR(30),
      created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE credits (
     id BIGSERIAL PRIMARY KEY,
     account_id BIGINT REFERENCES accounts(id),
     principal NUMERIC(15,2),
     rate NUMERIC(5,2),
     term_months INT,
     created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE payment_schedules (
   id BIGSERIAL PRIMARY KEY,
   credit_id BIGINT REFERENCES credits(id),
   due_date TIMESTAMP,
   amount NUMERIC(15,2),
   paid BOOLEAN DEFAULT FALSE
);