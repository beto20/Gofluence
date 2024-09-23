package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func UploadImage(container string, cs string, imagePath string) string {
	ctx := context.Background()
	client, err := azblob.NewClientFromConnectionString(cs, nil)

	handleError(err)

	// Specify the path to the PNG file
	filePath := imagePath

	// Open the PNG file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
	}
	defer file.Close()

	// Read the file contents into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	// Upload to data to blob storage
	_, err = client.UploadBuffer(ctx, container, filePath, fileBytes, &azblob.UploadBufferOptions{})
	handleError(err)

	remotePath := client.URL() + container + "/" + filePath
	fmt.Println("remotePath", remotePath)
	return remotePath
}
