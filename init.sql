CREATE
EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    ID       SERIAL NOT NULL PRIMARY KEY,
    nickname CITEXT NOT NULL UNIQUE COLLATE "POSIX",
    fullname TEXT   NOT NULL,
    email    CITEXT   NOT NULL UNIQUE,
    about    TEXT
);
--indexes
CREATE INDEX idx_nick_nick ON users (nickname);
CREATE INDEX idx_nick_email ON users (email);
CREATE INDEX idx_nick_cover ON users (nickname, fullname, about, email);

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
--indexes
CREATE INDEX idx_forum_slug ON forums using hash(slug);

DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE threads
(
    ID      SERIAL                          NOT NULL PRIMARY KEY,
    author  CITEXT                          NOT NULL REFERENCES users (nickname),
 --   created TEXT                            NOT NULL,
   created TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    forum   CITEXT REFERENCES forums (slug) NOT NULL,
    message TEXT                            NOT NULL,
    slug    CITEXT UNIQUE,
    title   TEXT                            NOT NULL,
    votes   INTEGER DEFAULT 0
);
--indexes
CREATE INDEX idx_thread_id ON threads(id);
CREATE INDEX idx_thread_slug ON threads(slug);
CREATE INDEX idx_thread_coverage ON threads (forum, created, id, slug, author, title, message, votes);

DROP TABLE IF EXISTS posts;
CREATE TABLE posts
(
    ID      SERIAL                          NOT NULL PRIMARY KEY,
    author  CITEXT                          NOT NULL REFERENCES users (nickname),
    created TIMESTAMP WITH TIME ZONE,
    edited  BOOLEAN DEFAULT false           NOT NULL,
    forum   CITEXT REFERENCES forums (slug) NOT NULL,
    message TEXT                            NOT NULL,
    parent  INTEGER DEFAULT 0               NOT NULL,
    thread  INTEGER REFERENCES threads (ID) NOT NULL,
    path    INTEGER[] DEFAULT '{0}':: INTEGER [] NOT NULL
);
--indexes
CREATE INDEX ON posts(thread, id, created, author, edited, message, parent, forum);
CREATE INDEX idx_post_thread_id_p_i ON posts (thread, (path[1]), id);

DROP TABLE IF EXISTS forum_users;
CREATE TABLE forum_users
(
    forumID INTEGER REFERENCES forums (ID),
    userID  INTEGER REFERENCES users (ID)
);
ALTER TABLE IF EXISTS forum_users ADD CONSTRAINT uniq UNIQUE (forumID, userID);
CREATE INDEX idx_forum_user ON forum_users (forumID, userID);

DROP TABLE IF EXISTS votes;
CREATE TABLE votes
(
    user_nick CITEXT REFERENCES users (nickname) NOT NULL,
    voice BOOLEAN NOT NULL,
    thread  INTEGER REFERENCES threads (ID) NOT NULL
  --  CONSTRAINT unique_vote UNIQUE (user_nick, thread)
);
ALTER TABLE IF EXISTS votes ADD CONSTRAINT uniq_votes UNIQUE (user_nick, thread);
CREATE INDEX idx_vote ON votes(thread, voice);

