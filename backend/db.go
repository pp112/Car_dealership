package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "os"
    "time"
)

func mustGetenv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}

func connectWithRetry() (*sql.DB, error) {
    host := mustGetenv("DB_HOST", "db")
    port := mustGetenv("DB_PORT", "5432")
    user := mustGetenv("DB_USER", "myuser")
    pass := mustGetenv("DB_PASSWORD", "mypassword")
    dbname := mustGetenv("DB_NAME", "myappdb")
    ssl := mustGetenv("DB_SSLMODE", "disable")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, pass, dbname, ssl)

    var db *sql.DB
    var err error
    for i := 0; i < 20; i++ {
        db, err = sql.Open("postgres", dsn)
        if err == nil {
            err = db.Ping()
        }
        if err == nil {
            return db, nil
        }
        time.Sleep(1 * time.Second)
    }
    return nil, err
}
