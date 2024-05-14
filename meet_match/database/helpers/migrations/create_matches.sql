CREATE TABLE IF NOT EXISTS matches(
    id SERIAL,
    session_id uuid,
    datetime timestamp NOT NULL,
    got_feedback BOOL,
    matched_card_id bigint,
    PRIMARY KEY (id)
);