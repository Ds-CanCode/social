CREATE TABLE folowers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user1 INTEGER NOT NULL,
    user2 INTEGER NOT NULL,
    accepted BOOLEAN NOT NULL DEFAULT 0,
    is_read INTEGER DEFAULT 0,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user1) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (user2) REFERENCES users (id) ON DELETE CASCADE
);