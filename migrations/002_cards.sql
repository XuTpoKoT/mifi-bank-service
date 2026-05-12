CREATE TABLE cards (
   id BIGSERIAL PRIMARY KEY,
   account_id BIGINT NOT NULL REFERENCES accounts(id),
   encrypted_pan TEXT NOT NULL,
   encrypted_expiry TEXT NOT NULL,
   cvv_hash TEXT NOT NULL,
   created_at TIMESTAMP DEFAULT now()
);