-- rambler up

CREATE TYPE stamp_state AS enum('pending', 'processing', 'sent', 'confirmed', 'failed');

CREATE TABLE IF NOT EXISTS stamps (
    id          UUID PRIMARY KEY,
    txhash      TEXT,
    merkletree  jsonb,
    state       stamp_state NOT NULL DEFAULT 'pending'
);

ALTER TABLE IF EXISTS stamp_requests
    ADD COLUMN IF NOT EXISTS stampId UUID;

-- rambler down

DROP TABLE IF EXISTS stamps;
DROP TYPE IF EXISTS stamp_state;
