CREATE TABLE IF NOT EXISTS matches(
    id SERIAL,
    session_id uuid,
    datetime timestamp NOT NULL,
    got_feedback BOOL,
    matched_card_id bigint,
    user_id bigint,
	match_viewed   bool,
    PRIMARY KEY (id)
);