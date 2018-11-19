-- rambler up

CREATE TYPE document_state AS enum('pending', 'batched', 'failed');

CREATE TABLE IF NOT EXISTS documents (
    id      UUID PRIMARY KEY,
    hash    TEXT,
    state   document_state
);

-- rambler down

DROP TABLE IF EXISTS document_state;
DROP TYPE IF EXISTS document_state;
