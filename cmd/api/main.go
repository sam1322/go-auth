package main

import (
	"fmt"
	"goAuth/internal/auth"
	"goAuth/internal/server"
)

func main() {

	auth.NewAuth()
	// storage.Example_client_NewClient()
	// storage.Example_client_NewClientWithSharedKeyCredential()
	// storage.Example_client_NewClientFromConnectionString()
	// storage.GetFileShareClient()
	// storage.AzureFileUpload()
	fmt.Println("Starting server...")

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
