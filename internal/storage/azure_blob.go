package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func handleError(err error) {
	if err != nil {
		// log.Fatal(err.Error())
		log.Println(err.Error())
	}
}

func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", filename, line, err)
		b = true
	}
	return
}

// this logs the function name as well.
func FancyHandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, filename, line, _ := runtime.Caller(1)

		log.Printf("[error] in %s \n[%s:%d]\n Error %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		b = true
	}
	return
}

func Example_client_NewClient() {
	// this example uses Azure Active Directory (AAD) to authenticate with Azure Blob Storage
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DefaultAzureCredential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fmt.Println(cred, *cred, &cred)

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// fmt.Println

	fmt.Println(client.URL())
}

func Example_client_NewClientWithSharedKeyCredential() {
	// this example uses a shared key to authenticate with Azure Blob Storage
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	// accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// shared key authentication requires the storage account name and access key
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	serviceClient, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func blobClientFromConnectionString() {
	// this example uses a connection string to authenticate with Azure Blob Storage
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func getBlobServiceClient() (*azblob.Client, error) {
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return &azblob.Client{}, err
	}
	return serviceClient, nil
}

func Example_client_UploadFile() {
	// Set up file to upload
	fileSize := 8 * 1024 * 1024
	fileName := "test_upload_file2.txt"
	fileData := make([]byte, fileSize)
	fileData = []byte("hello\ngo\n")
	// err := os.WriteFile(fileName, fileData, 0666)
	err := os.WriteFile(fileName, fileData, 0644)
	handleError(err)

	// Open the file to upload
	fileHandler, err := os.Open(fileName)
	handleError(err)

	// close the file after it is no longer required.
	defer func(file *os.File) {
		err = file.Close()
		handleError(err)
	}(fileHandler)

	// delete the local file if required.
	defer func(name string) {
		err = os.Remove(name)
		handleError(err)
	}(fileName)

	// return

	// using azure identity to authenticate with Azure Blob Storage

	// accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	// if !ok {
	// 	panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	// }
	// serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// cred, err := azidentity.NewDefaultAzureCredential(nil)
	// handleError(err)

	// client, err := azblob.NewClient(serviceURL, cred, nil)
	// handleError(err)

	// using a shared key to authenticate with Azure Blob Storage

	// accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	// if !ok {
	// 	panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	// }
	// // accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	// accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	// if !ok {
	// 	panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	// }
	// serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// // shared key authentication requires the storage account name and access key
	// cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	// handleError(err)
	// client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)

	client, err := getBlobServiceClient()
	handleError(err)

	// Upload the file to a block blob
	_, err = client.UploadFile(context.TODO(), "sam-blob", fileName, fileHandler,
		// _, err = client.UploadFile(context.TODO(), "demo-share", fileName, fileHandler,
		// _, err = client.UploadFile(context.TODO(), "sam-blob", "directory/"+fileName, fileHandler,
		&azblob.UploadFileOptions{
			BlockSize:   int64(1024),
			Concurrency: uint16(3),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Println(bytesTransferred)
			},
		})
	handleError(err)
}

func BlobUploadMultipartFile() error {
	client, err := getBlobServiceClient()
	// handleError(err)
	if FancyHandleError(err) {
		log.Print("stuff")
		return err
	}
	containerName := "sam-blob"
	blobData := "Hello world!"
	blobName := "HelloWorld.txt"

	now := time.Now()
	timestamp := now.Format("2006-01-02T15-04-05") // Adjust format as needed

	blobName = fmt.Sprintf("file-%v.txt", now.UnixMilli())
	blobName = fmt.Sprintf("%s-%s.txt", "file", timestamp)
	uploadResp, err := client.UploadStream(context.TODO(),
		containerName,
		blobName,
		strings.NewReader(blobData),
		&azblob.UploadStreamOptions{
			Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
			Tags:     map[string]string{"Year": "2024"},
		})
	if FancyHandleError(err) {
		log.Print("stuff")
		return err
	}
	fmt.Println(uploadResp)
	return nil
}

// // This example is a quick-starter and demonstrates how to get started using the Azure Blob Storage SDK for Go.
// func Example() {
// 	// Your account name and key can be obtained from the Azure Portal.
// 	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
// 	if !ok {
// 		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
// 	}

// 	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
// 	if !ok {
// 		panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
// 	}
// 	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
// 	handleError(err)

// 	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
// 	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
// 	handleError(err)

// 	// fmt.Println(&client, client, *client)
// 	return
// 	// ===== 1. Create a container =====
// 	containerName := "testcontainer"
// 	containerCreateResp, err := client.CreateContainer(context.TODO(), containerName, nil)
// 	handleError(err)
// 	fmt.Println(containerCreateResp)

// 	// ===== 2. Upload and Download a block blob =====
// 	blobData := "Hello world!"
// 	blobName := "HelloWorld.txt"
// 	uploadResp, err := client.UploadStream(context.TODO(),
// 		containerName,
// 		blobName,
// 		strings.NewReader(blobData),
// 		&azblob.UploadStreamOptions{
// 			Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
// 			Tags:     map[string]string{"Year": "2022"},
// 		})
// 	handleError(err)
// 	fmt.Println(uploadResp)

// 	// Download the blob's contents and ensure that the download worked properly
// 	blobDownloadResponse, err := client.DownloadStream(context.TODO(), containerName, blobName, nil)
// 	handleError(err)

// 	// Use the bytes.Buffer object to read the downloaded data.
// 	// RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
// 	reader := blobDownloadResponse.Body
// 	downloadData, err := io.ReadAll(reader)
// 	handleError(err)
// 	if string(downloadData) != blobData {
// 		log.Fatal("Uploaded data should be same as downloaded data")
// 	}

// 	err = reader.Close()
// 	if err != nil {
// 		return
// 	}

// 	// ===== 3. List blobs =====
// 	// List methods returns a pager object which can be used to iterate over the results of a paging operation.
// 	// To iterate over a page use the NextPage(context.Context) to fetch the next page of results.
// 	// PageResponse() can be used to iterate over the results of the specific page.
// 	pager := client.NewListBlobsFlatPager(containerName, nil)
// 	for pager.More() {
// 		resp, err := pager.NextPage(context.TODO())
// 		handleError(err)
// 		for _, v := range resp.Segment.BlobItems {
// 			fmt.Println(*v.Name)
// 		}
// 	}

// 	// Delete the blob.
// 	_, err = client.DeleteBlob(context.TODO(), containerName, blobName, nil)
// 	handleError(err)

// 	// Delete the container.
// 	_, err = client.DeleteContainer(context.TODO(), containerName, nil)
// 	handleError(err)
// }

const random64BString string = "2SDgZj6RkKYzJpu04sweQek4uWHO8ndPnYlZ0tnFS61hjnFZ5IkvIGGY44eKABov"
