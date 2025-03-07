CREATE TABLE
    EventStatus (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        member_id INTEGER NOT NULL,
        event_id INTEGER NOT NULL,
        group_id INTEGER NOT NULL,
        FOREIGN KEY (member_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (event_id) REFERENCES posts_event (id) ON DELETE CASCADE ON UPDATE CASCADE
    );