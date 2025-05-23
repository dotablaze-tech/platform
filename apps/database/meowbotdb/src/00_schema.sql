-- CREATE DATABASE meowbot;

CREATE TABLE guilds
(
    id         TEXT PRIMARY KEY,
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
    guild_id            TEXT REFERENCES guilds (id) ON DELETE CASCADE,
    user_id             TEXT REFERENCES users (id) ON DELETE CASCADE,
    successful_meows    INT DEFAULT 0,
    failed_meows        INT DEFAULT 0,
    total_meows         INT DEFAULT 0,
    current_streak      INT DEFAULT 0,
    highest_streak      INT DEFAULT 0,
    last_meow_at        TIMESTAMP,
    last_failed_meow_at TIMESTAMP,
    PRIMARY KEY (guild_id, user_id)
);

CREATE INDEX idx_user_guild_stats_user_id ON user_guild_stats (user_id);
CREATE INDEX idx_user_guild_stats_guild_id ON user_guild_stats (guild_id);

CREATE TABLE guild_streaks
(
    guild_id           TEXT PRIMARY KEY REFERENCES guilds (id) ON DELETE CASCADE,
    meow_count         INT DEFAULT 0,
    last_user_id       TEXT,
    high_score         INT DEFAULT 0,
    high_score_user_id TEXT
);

CREATE TABLE IF NOT EXISTS guild_channels
(
    guild_id   TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL
);

