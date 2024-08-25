-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (
    id VARCHAR PRIMARY KEY,
    tgid INTEGER,
    name VARCHAR,
    tgusername VARCHAR,
    chatid VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR,

    UNIQUE(tgid)
);
CREATE TABLE IF NOT EXISTS event (
    id VARCHAR PRIMARY KEY,
    chatid VARCHAR,
    ownerid INTEGER,
    text VARCHAR,
    notifyat VARCHAR,
    delta VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR
);
CREATE TABLE IF NOT EXISTS schedule (
    id VARCHAR PRIMARY KEY,
    chatid VARCHAR,
    ownerid INTEGER,
    text VARCHAR,
    delta VARCHAR,
    day VARCHAR,
    eventId VARCHAR,
    createdat VARCHAR,
    updatedat VARCHAR
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
