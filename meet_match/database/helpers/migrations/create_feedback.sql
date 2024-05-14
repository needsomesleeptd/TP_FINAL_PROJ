CREATE TABLE IF NOT EXISTS feedbacks(
    id SERIAL,
    user_id INT NOT NULL,
    description TEXT,
    hasgone BOOL NOT NULL,
    datetime timestamp NOT NULL,
    PRIMARY KEY (id)
    
);