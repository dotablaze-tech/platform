-- CREATE DATABASE meowbot;

CREATE TABLE guilds
(
    id         TEXT PRIMARY KEY,
    name       TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users
(
    id         TEXT PRIMARY KEY,
    username   TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_guild_stats
(
    guild_id         TEXT REFERENCES guilds (id),
    user_id          TEXT REFERENCES users (id),
    successful_meows INT DEFAULT 0,
    failed_meows     INT DEFAULT 0,
    highest_streak   INT DEFAULT 0,
    last_meow_at     TIMESTAMP,
    PRIMARY KEY (guild_id, user_id)
);
