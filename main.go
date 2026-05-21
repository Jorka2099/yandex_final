package main

import (
	"log"
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)



func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("Failed to initialize database", err)
	}

	server.Run()

	

}
