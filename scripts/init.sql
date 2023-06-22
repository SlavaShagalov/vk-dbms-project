CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id       bigserial PRIMARY KEY,
    nickname citext NOT NULL UNIQUE,
    fullname text   NOT NULL,
    about    text,
    email    citext NOT NULL UNIQUE
);
