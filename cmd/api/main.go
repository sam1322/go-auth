package main

import (
	"fmt"
	"goAuth/internal/auth"
	"goAuth/internal/server"
	"log"
)

func main() {

	auth.NewAuth()
	// storage.Example_client_NewClient()
	// storage.Example_client_NewClientWithSharedKeyCredential()
	// storage.Example_client_NewClientFromConnectionString()
	// storage.GetFileShareClient()
	// storage.AzureFileUpload()
	fmt.Println("Starting Server at 8080...")

	server := server.InitNewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
	log.Println("Server started successfully!")
}
