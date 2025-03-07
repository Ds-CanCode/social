package funcs

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/sqlite3"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const dbPath = "internal/database/tables.db"

func SetupDb() error {
    var err error
    DB, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatalf("Erreur connexion DB: %v", err)
    }
    log.Println("Connexion à SQLite réussie !")
    
    if _, err = DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
        return err
    }
    return ApplyMigrations()
}

func GetDb() *sql.DB {
    return DB
}

func ApplyMigrations() error {
    wd, err := os.Getwd()
    if err != nil {
        return fmt.Errorf(" Impossible d'obtenir le répertoire actuel: %w", err)
    }
    migrationsPath := "file://" + filepath.Join(wd, "migrations")

    m, err := migrate.New(migrationsPath, "sqlite3://"+dbPath)
    if err != nil {
        return fmt.Errorf(" Erreur chargement migrations: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf(" Erreur application migrations: %w", err)
    }

    log.Println("✅ Migrations appliquées avec succès !")
    return nil
}