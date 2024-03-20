CREATE TABLE IF NOT EXISTS users (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Balance FLOAT4 DEFAULT 0.0
);

CREATE TABLE IF NOT EXISTS quests (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Cost FLOAT4 DEFAULT 0.0
);

CREATE TABLE IF NOT EXISTS completedTasks (
    Id SERIAL PRIMARY KEY,
    userId integer not null references users(Id),
    questId integer not null references quests(Id)
);