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
     user_id BIGINT NOT NULL REFERENCES users(id),
     account_id BIGINT NOT NULL REFERENCES accounts(id),
     principal NUMERIC(15,2) NOT NULL,
     annual_rate NUMERIC(5,2) NOT NULL,
     term_months INT NOT NULL,
     monthly_payment NUMERIC(15,2) NOT NULL,
     remaining_debt NUMERIC(15,2) NOT NULL,
     status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
     created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE payment_schedules (
   id BIGSERIAL PRIMARY KEY,
   credit_id BIGINT NOT NULL REFERENCES credits(id),
   due_date DATE NOT NULL,
   amount NUMERIC(15,2) NOT NULL,
   status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
   penalty NUMERIC(15,2) NOT NULL DEFAULT 0,
   paid_at TIMESTAMP
);