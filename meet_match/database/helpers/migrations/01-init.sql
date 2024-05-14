CREATE TABLE IF NOT EXISTS places(
    place_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, -- generated
    price TEXT, -- event
    place INT, -- event
    url VARCHAR(254) NOT NULL, -- event, place
    title TEXT NOT NULL, -- event, place
    description TEXT, -- event, place
    subway VARCHAR(400), -- place
    timetable VARCHAR(200), -- place
    dates TEXT, -- event
    age_restriction VARCHAR(30), -- event
    categories TEXT, -- event, place
    phone VARCHAR(100), -- place
    foreign_url VARCHAR(254), -- event, place (events site_url == place foreign_url )
    favorites_count INT -- event, place 
);


CREATE TABLE IF NOT EXISTS embeddings(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    place_id INT NOT NULL REFERENCES places(place_id) ON DELETE CASCADE,
    embedding bytea NOT NULL
);


CREATE TABLE IF NOT EXISTS fact_scrolled(
    session_id VARCHAR(254) NOT NULL,
    user_id INT NOT NULL,
    place_id INT NOT NULL,
    is_liked boolean NOT NULL,
    datetime timestamp NOT NULL,
    PRIMARY KEY (session_id, user_id, place_id)
);