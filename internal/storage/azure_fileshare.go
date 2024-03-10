package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

func SanitizeFilename(filename string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9_.-]+") // Match characters other than letters, numbers, ., -, _
	return re.ReplaceAllString(filename, "-")    // Replace with "-"
}

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
	handleError(err)

	// shareName := "demo-share"

	now := time.Now()

	srcFileName := fmt.Sprintf("file-%v.txt", now.UnixMilli())

	fmt.Println("filename", srcFileName)
	// return
	// fileSize := int64(5)

	handleError(err)

	fileData := []byte("hello\ngo\n")
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
	fileSize := fInfo.Size()

	defer func(fh *os.File) {
		err := fh.Close()
		handleError(err)
	}(fh)

	dirName := "temp"
	// dirClient := shareClient.NewRootDirectoryClient() //create or get the root directory
	dirClient := shareClient.NewDirectoryClient(dirName)
	_, err = dirClient.Create(context.TODO(), nil) // to create directory if it does not exists

	srcFileClient := dirClient.NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil) // to create file if it does not exists

	// handleError(err)
	if FancyHandleError(err) {
		log.Print("stuff")
		return
	}

	err = srcFileClient.UploadFile(context.Background(), fh, nil)
	if err != nil {
		fmt.Println(err)
		log.Println("failed to upload file")
	}
	fmt.Println("File Upload Successfuly", srcFileClient.URL())

}

func AzureFileUploadByBytes(fileBuffer []byte, fileSize int64, fileName string) error {
	shareClient, err := getFileShareClient()
	handleError(err)

	// shareName := "demo-share"

	// Generate a timestamp string
	now := time.Now()
	timestamp := now.Format("2006-01-02T15-04-05") // Adjust format as needed

	srcFileName := fmt.Sprintf("file-%v.txt", now.UnixMilli())

	srcFileName = fileName
	srcFileName = fmt.Sprintf("%s-%s", timestamp, fileName)
	srcFileName = SanitizeFilename(srcFileName)

	fmt.Println("filename", srcFileName, "fileSize", fileSize, "FileSize in MB", fileSize/1024/1024, "MB")
	handleError(err)

	dirName := "temp"

	// dirClient := shareClient.NewRootDirectoryClient()
	dirClient := shareClient.NewDirectoryClient(dirName)
	_, err = dirClient.Create(context.TODO(), nil) // to create directory if it does not exists

	srcFileClient := dirClient.NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil) // to create file if it does not exists
	fmt.Println("File Created", srcFileClient.URL())
	if FancyHandleError(err) {
		log.Print("stuff")
		return err
	}

	err = srcFileClient.UploadBuffer(context.Background(), fileBuffer, nil)

	// handleError(err)
	if FancyHandleError(err) {
		log.Print("stuff")
		return err
	}

	if err != nil {
		fmt.Println(err)
		log.Println("failed to upload file")
	}
	fmt.Println("File Upload Successfuly", srcFileClient.URL())
	return nil
	// _, err = srcFileClient.Delete(context.Background(), nil) // to delete the file if we want to
	// handleError(err)

	// _, err = shareClient.Delete(context.Background(), nil) // to delete the directory if we want to but it might give an error about snapshot
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

func GetSASUrl() (string, error) {
	shareClient, err := getFileShareClient()
	if FancyHandleError(err) {
		log.Print("stuff")
		return "", err
	}
	dirName := "temp"
	srcFileName := "2024-03-10T13-39-39-trafagar-d-law.jpg"

	dirClient := shareClient.NewDirectoryClient(dirName)
	_, err = dirClient.Create(context.TODO(), nil) // to create directory if it does not exists

	srcFileClient := dirClient.NewFileClient(srcFileName)
	// _, err = srcFileClient.Create(context.Background(), fileSize, nil) // to create file if it does not exists
	permission := sas.FilePermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(1, 0, 0)
	options := file.GetSASURLOptions{StartTime: &start}
	sasURL, err := srcFileClient.GetSASURL(permission, expiry, &options)
	if FancyHandleError(err) {
		log.Print("stuff")
		return "", err
	}
	// _ = sasURL
	fmt.Println("SAS URL", sasURL)
	return sasURL, nil
}
