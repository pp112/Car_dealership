package main

import (
    "log"
    "os"
)

func main() {
    db, err := connectWithRetry()
    if err != nil {
        log.Fatalf("failed to connect to db: %v", err)
    }
    defer db.Close()

    srv := &Server{DB: db}
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    r := srv.routes()
    log.Printf("starting server on :%s ...", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
