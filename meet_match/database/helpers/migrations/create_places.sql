-- CREATE TABLE IF NOT EXISTS places(
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     url VARCHAR(254) NOT NULL,
--     title TEXT NOT NULL
-- );

-- events site_url == foreign_url 
CREATE TABLE IF NOT EXISTS places(
    place_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, -- generated
    price INT, -- event
    place INT, -- event
    url VARCHAR(254) NOT NULL, -- event, place
    title TEXT NOT NULL, -- event, place
    description TEXT, -- event, place
    subway VARCHAR(400), -- place
    timetable VARCHAR(200), -- place
    dates TEXT, -- event
    age_restriction VARCHAR(30), -- event
    categories TEXT[], -- event, place
    phone VARCHAR(100), -- place
    foreign_url VARCHAR(254), -- event, place (events site_url == place foreign_url )
    favorites_count INT -- event, place 
);


CREATE TABLE IF NOT EXISTS embeddings(
    id INT GENERATED ALWAYS AS IDENTITY
    place_id INT NOT NULL REFERENCES places(place_id) ON DELETE CASCADE,
    embedding bytea NOT NULL
);
