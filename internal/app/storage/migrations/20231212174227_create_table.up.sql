CREATE TABLE links (
    id serial PRIMARY KEY,
    long varchar(1000) NOT NULL,
    short varchar(12) NOT NULL UNIQUE
)