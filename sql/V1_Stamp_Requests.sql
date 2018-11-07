-- rambler up

CREATE TYPE stamp_request_state AS enum('pending', 'batched', 'failed');

CREATE TABLE IF NOT EXISTS stamp_requests (
    id UUID PRIMARY KEY,
    hash TEXT,
    state stamp_request_state
);

-- rambler down

DROP TABLE IF EXISTS stamp_requests;
DROP TYPE IF EXISTS stamp_request_state;
