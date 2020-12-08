CREATE
EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    ID       SERIAL NOT NULL PRIMARY KEY,
    nickname CITEXT NOT NULL UNIQUE,
    fullname TEXT   NOT NULL,
    email    TEXT   NOT NULL UNIQUE,
    about    TEXT
);

DROP TABLE IF EXISTS forums CASCADE;

CREATE TABLE forums
(
    ID        SERIAL                             NOT NULL PRIMARY KEY,
    slug      CITEXT                             NOT NULL UNIQUE,
    threads   INTEGER DEFAULT 0                  NOT NULL,
    posts     INTEGER DEFAULT 0                  NOT NULL,
    title     TEXT                               NOT NULL,
    user_nick CITEXT REFERENCES users (nickname) NOT NULL
);

DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE threads
(
    ID      SERIAL                                 NOT NULL PRIMARY KEY,
    author  CITEXT                                 NOT NULL REFERENCES users (nickname),
    created TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    forum   CITEXT REFERENCES forums (slug)        NOT NULL,
    message TEXT                                   NOT NULL,
    slug    CITEXT UNIQUE,
    title   TEXT                                   NOT NULL,
    votes   INTEGER                  DEFAULT 0
);

DROP TABLE IF EXISTS posts;
CREATE TABLE posts
(
    ID      SERIAL                          NOT NULL,
    author  CITEXT                          NOT NULL REFERENCES users (nickname),
    created TIMESTAMP WITH TIME ZONE,
    edited  BOOLEAN DEFAULT false           NOT NULL,
    forum   CITEXT REFERENCES forums (slug) NOT NULL,
    message TEXT                            NOT NULL,
    parent  INTEGER DEFAULT 0               NOT NULL,
    thread  INTEGER REFERENCES threads (ID) NOT NULL,
    path    INTEGER[] DEFAULT '{0}':: INTEGER [] NOT NULL
);

DROP TABLE IF EXISTS forum_users;
CREATE TABLE forum_users
(
    forumID INTEGER REFERENCES forums (ID),
    userID  INTEGER REFERENCES users (ID)
);
ALTER TABLE IF EXISTS forum_users ADD CONSTRAINT uniq UNIQUE (forumID, userID)

DROP TABLE IF EXISTS votes;
CREATE TABLE votes
(
    user_nick CITEXT REFERENCES users (nickname) NOT NULL,
    voice BOOLEAN NOT NULL,
    thread  INTEGER REFERENCES threads (ID) NOT NULL,
    CONSTRAINT unique_vote UNIQUE (user_nick, thread)
);