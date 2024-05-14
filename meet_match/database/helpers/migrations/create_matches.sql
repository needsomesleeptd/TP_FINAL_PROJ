CREATE TABLE IF NOT EXISTS matches(
    id SERIAL,
    session_id uuid,
    datetime timestamp NOT NULL,
    got_feedback BOOL,
    card_matched_id uint64,
    PRIMARY KEY (id),
);