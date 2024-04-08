CREATE TABLE IF NOT EXISTS fact_scrolled(
    session_id VARCHAR(254) NOT NULL,
    user_id INT NOT NULL,
    places_id INT NOT NULL,
    is_liked boolean NOT NULL,
    PRIMARY KEY (session_id, user_id, places_id)
);