CREATE TABLE
    groups_messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        group_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        createdAt DATETIME NOT NULL,
        FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE
    );