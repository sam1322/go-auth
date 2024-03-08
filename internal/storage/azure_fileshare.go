package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

func GetFileShareClient() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "demo-share"
	filePath := "download.jpeg"
	fileClient, err := file.NewClientFromConnectionString(connectionString, shareName, filePath, nil)
	handleError(err)
	fmt.Println(fileClient.URL())
}

func getFileShareClient() (*share.Client, error) {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "demo-share"
	// filePath := "download.jpeg"
	// fileClient, err := share.NewClientFromConnectionString(connectionString, shareName, filePath, nil)
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)

	// handleError(err)
	if err != nil {
		return &share.Client{}, err
	}
	fmt.Println(shareClient.URL())
	return shareClient, nil
}

func AzureFileUpload() {
	shareClient, err := getFileShareClient()
	// handleError(err)

	// shareName := "demo-share"
	srcFileName := "download3.txt"
	fileSize := int64(5)

	// connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	// if !ok {
	// 	log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	// }

	// shareClient, err = share.NewClientFromConnectionString(connectionString, shareName, nil)
	// handleError(err)

	// _, err = shareClient.Create(context.Background(), nil)
	// fmt.Println("1")
	// handleError(err)

	// dirName := "temp"

	// dirClient := shareClient.NewDirectoryClient(dirName)
	// _, err = dirClient.Create(context.TODO(), nil)
	// fmt.Println("3")

	handleError(err)

	// _, content := generateData(int(fileSize))
	fileData := make([]byte, fileSize)
	fileData = []byte("hello\ngo\n")
	content := fileData
	err = os.WriteFile(srcFileName, content, 0644)
	handleError(err)

	defer func() {
		err = os.Remove(srcFileName)
		handleError(err)
	}()

	fh, err := os.Open(srcFileName)

	handleError(err)
	// get the size of file
	fInfo, err := fh.Stat()
	// TODO: handle error
	handleError(err)
	fileSize = fInfo.Size()

	defer func(fh *os.File) {
		err := fh.Close()
		handleError(err)
	}(fh)
	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)

	fmt.Println("2")

	handleError(err)

	err = srcFileClient.UploadFile(context.Background(), fh, nil)
	if err != nil {
		fmt.Println(err)
		log.Println("failed to upload file")
	}
	fmt.Println("File Upload Successfuly")

	// err = srcFileClient.UploadFile(context.Background(), fh, nil)

	// Create a file
	// fileContents := []byte("Hello Azure!")
	// _, err := fileClient.Upload(context.Background(), bytes.NewReader(fileContents), true)
	// handleError(err)
}

func generateData(sizeInBytes int) (io.ReadSeekCloser, []byte) {
	data := make([]byte, sizeInBytes)
	_len := len(random64BString)
	if sizeInBytes > _len {
		count := sizeInBytes / _len
		if sizeInBytes%_len != 0 {
			count = count + 1
		}
		copy(data[:], strings.Repeat(random64BString, count))
	} else {
		copy(data[:], random64BString)
	}
	return streaming.NopCloser(bytes.NewReader(data)), data
}
