package main

import (
	"log"

	"tradingsystem/internal/database"
	"tradingsystem/internal/router"
)

func main() {
	db, err := database.OpenMySQLFromEnv()
	if err != nil {
		log.Fatalf("open database failed: %v", err)
	}

	engine := router.New(db)

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
