-- CREATE TABLE IF NOT EXISTS places(
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     url VARCHAR(254) NOT NULL,
--     title TEXT NOT NULL
-- );

-- id,title,short_title
CREATE TABLE IF NOT EXISTS places(
    place_id INT NOT NULL PRIMARY KEY,
    url VARCHAR(254) NOT NULL,
    title TEXT NOT NULL,
    description TEXT
);


CREATE TABLE IF NOT EXISTS embeddings(
                                         place_id INT NOT NULL PRIMARY KEY,
                                         embedding bytea NOT NULL
);
