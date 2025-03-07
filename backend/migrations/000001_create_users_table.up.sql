CREATE TABLE
    users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email VARCHAR(255) NOT NULL UNIQUE,
        password CHAR(60) NOT NULL,
        firstName VARCHAR(255) NOT NULL,
        lastName VARCHAR(255) NOT NULL,
        datebirth DATE NOT NULL,
        avatar TEXT DEFAULT '',
        nickname VARCHAR(255) DEFAULT '',
        aboutme VARCHAR(255) DEFAULT '',
        profileType BOOLEAN DEFAULT 0,
        createdAt DATETIME NOT NULL
    );