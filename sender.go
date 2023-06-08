package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	chunkSize = 1024 * 1024 // 1 MB
)

func main() {
	filePaths := []string{
		`C:\Users\Admin\Downloads\_import_609116af6c46b1.75389671.mov`,
		`C:\Users\Admin\Downloads\_import_62fdc893387585.82713010.mov`,
		`C:\Users\Admin\Downloads\221206_02_Currency_4k_002.mp4`,
	}

	for _, filePath := range filePaths {
		err := processFile(filePath)
		if err != nil {
			log.Println("Error processing file:", err)
		}
	}
}

func processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file information: %v", err)
	}

	fileName := fileInfo.Name()
	outputFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outputFile.Close()

	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("failed to read file: %v", err)
			}
			break
		}
		if n > 0 {
			// Send chunk to receiver microservice
			err := sendChunkToReceiver(buf[:n])
			if err != nil {
				return fmt.Errorf("failed to send chunk to receiver: %v", err)
			}
		}
	}

	fmt.Printf("File %s processed successfully\n", fileName)
	return nil
}

func sendChunkToReceiver(chunk []byte) error {
	// Implement the logic to send the chunk to the receiver microservice
	// You can use HTTP, gRPC, or any other communication protocol of your choice
	// Make sure the receiver microservice has an appropriate endpoint to receive the chunk
	return nil
}
