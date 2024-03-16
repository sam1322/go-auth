package main

import (
	"fmt"
	"goAuth/internal/auth"
	"goAuth/internal/server"
	"log"
)

func main() {

	auth.NewAuth()
	fmt.Println("Starting Server at 8080...")
	// todo.AddTodo()

	server := server.InitNewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
	log.Println("Server started successfully!")
}
