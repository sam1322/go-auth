package main

import (
	"fmt"
	"goAuth/internal/auth"
	"goAuth/internal/server"
	"goAuth/internal/storage"
)

func main() {

	auth.NewAuth()
	// storage.Example_client_NewClient()
	// storage.Example_client_NewClientWithSharedKeyCredential()
	// storage.Example_client_NewClientFromConnectionString()
	// storage.GetFileShareClient()
	storage.AzureFileUpload()

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
