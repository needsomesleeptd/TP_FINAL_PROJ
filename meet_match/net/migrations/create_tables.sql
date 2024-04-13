CREATE TABLE IF NOT EXISTS fact_scrolled(
    id INT NOT NULL PRIMARY KEY,
    session_id VARCHAR(254) NOT NULL,
    user_id INT NOT NULL,
    places_id INT NOT NULL,
    is_liked boolean NOT NULL,
    datetime timestamp NOT NULL
);