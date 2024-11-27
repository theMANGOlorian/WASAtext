CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(25) NOT NULL UNIQUE,
    photo VARCHAR(27) UNIQUE,
    CHECK (LENGTH(username) >= 3)
);

CREATE TABLE IF NOT EXISTS conversations (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(7) NOT NULL,
    name VARCHAR(25),
    photo INTEGER,
    CHECK (( type = 'private' AND name IS NULL AND photo IS NULL) OR 
        (type = 'group' AND name IS NOT NULL AND photo IS NOT NULL )),
    CHECK ( type IN ('private','group'))
);

CREATE TABLE IF NOT EXISTS members (
    userId VARCHAR(36) NOT NULL,
    conversationId VARCHAR(36) NOT NULL,
    PRIMARY KEY (userId,conversationId),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id)
);

CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(5) NOT NULL,
    text VARCHAR(5000),
    photo INTEGER,
    conversation VARCHAR(36) NOT NULL,
    replay VARCHAR(36),
    status VARCHAR(4) NOT NULL,
    timestamp TEXT NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (conversation) REFERENCES conversations(id),
    FOREIGN KEY (replay) REFERENCES messages(id),
    CHECK ( status IN ('read','recv','none') ),
    CHECK ( type IN ('photo','text') ),
    CHECK ( (type = 'text' AND text IS NOT NULL AND photo IS NULL) OR 
        ( type = 'photo' AND photo IS NOT NULL AND text IS NULL))

);

CREATE TABLE IF NOT EXISTS users_msg (
    userId VARCHAR(36) NOT NULL,
    msgId VARCHAR(36) NOT NULL,
    status VARCHAR(4) NOT NULL,
    PRIMARY KEY (userId, msgId),
    CHECK (status IN ('read','recv','none')),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (msgId) REFERENCES messages(id)
);

CREATE TABLE IF NOT EXISTS reactions (
    id VARCHAR(36) PRIMARY KEY,
    owner VARCHAR(36) NOT NULL,
    messageId INTEGER NOT NULL,
    emoji TEXT NOT NULL,
    CHECK (emoji IN ('😊','😂','😍','😎','🥺')),
    UNIQUE(owner, messageId)
);