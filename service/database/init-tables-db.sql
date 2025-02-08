CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(25) NOT NULL UNIQUE,
    photo VARCHAR(27) UNIQUE,
    CHECK (LENGTH(username) >= 3)
);

CREATE TABLE IF NOT EXISTS conversations (
    ------------both (group & private)--------------
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(7) NOT NULL,
    --------------- group only ---------------------
    name VARCHAR(25),
    photo VARCHAR(27),
    --------------- private only -------------------
    -- none
    CHECK (( type = 'private' AND name IS NULL AND photo IS NULL) OR 
        (type = 'group' AND name IS NOT NULL)),
    CHECK ( type IN ('private','group'))
);

CREATE TABLE IF NOT EXISTS members (
    userId VARCHAR(36) NOT NULL,
    conversationId VARCHAR(36) NOT NULL,
    joinDate TEXT DEFAULT (datetime('now','localtime')) NOT NULL,
    tsLastRecv TEXT DEFAULT '0000-01-01 00:00:00' NOT NULL,
    tsLastRead TEXT DEFAULT '0000-01-01 00:00:00' NOT NULL,
    PRIMARY KEY (userId,conversationId),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id)
);

CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(36) PRIMARY KEY,
    sender VARCHAR(36) NOT NULL,
    type VARCHAR(5) NOT NULL,
    text VARCHAR(5000),
    photo VARCHAR(27),
    conversation VARCHAR(36) NOT NULL,
    reply VARCHAR(36),
    forwarded INTEGER DEFAULT 0 NOT NULL,
    status VARCHAR(4) DEFAULT 'none' NOT NULL,
    timestamp TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (conversation) REFERENCES conversations(id),
    FOREIGN KEY (reply) REFERENCES messages(id),
    CHECK ( status IN ('read','recv','none') ),
    CHECK (forwarded IN (0,1)),
    CHECK ( type IN ('photo','text') ),
    CHECK ( (type = 'text' AND text IS NOT NULL AND photo IS NULL) OR 
        ( type = 'photo' AND photo IS NOT NULL AND text IS NULL))
);

CREATE TABLE IF NOT EXISTS reactions (
    id VARCHAR(36) PRIMARY KEY,
    owner VARCHAR(36) NOT NULL,
    messageId INTEGER NOT NULL,
    emoji TEXT NOT NULL,
    CHECK (emoji IN ('😊','😂','😍','😎','😭')),
    UNIQUE(owner, messageId)
);

