CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id       bigserial,
    nickname citext NOT NULL UNIQUE PRIMARY KEY,
    fullname text   NOT NULL,
    about    text,
    email    citext NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forums
(
    id            bigserial,
    title         text   NOT NULL,
    slug          citext NOT NULL UNIQUE PRIMARY KEY,
    user_nickname citext
        CONSTRAINT user_nickname NOT NULL REFERENCES users (nickname),
    threads       int DEFAULT 0,
    posts         int DEFAULT 0
);

CREATE TABLE IF NOT EXISTS forum_users
(
    nickname citext NOT NULL COLLATE "ucs_basic" REFERENCES users (nickname),
    fullname text   NOT NULL,
    about    text   NOT NULL,
    email    citext NOT NULL,
    forum    citext NOT NULL REFERENCES forums (slug),
    PRIMARY KEY (nickname, forum)
);

CREATE TABLE IF NOT EXISTS threads
(
    id         bigserial PRIMARY KEY,
    title      text   NOT NULL,
    author     citext NOT NULL REFERENCES users (nickname),
    forum      citext NOT NULL REFERENCES forums (slug),
    message    text   NOT NULL,
    votes      int                      DEFAULT 0,
    slug       citext,
    created_at timestamp with time zone DEFAULT now()
);
