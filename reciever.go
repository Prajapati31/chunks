package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	chunkSize = 1024 * 1024 // 1 MB
)

func main() {
	http.HandleFunc("/receive", receiveHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	chunk, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := r.Header.Get("FileName")
	err = saveChunk(fileName, chunk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Chunk received successfully")
}

func saveChunk(fileName string, chunk []byte) error {
	// Create a directory to store the chunks for each file
	tempDir := "temp"
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Save the chunk to a file in the temp directory
	chunkPath := filepath.Join(tempDir, fileName)
	f, err := os.OpenFile(chunkPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to save chunk: %v", err)
	}
	defer f.Close()

	_, err = f.Write(chunk)
	if err != nil {
		return fmt.Errorf("failed to write chunk: %v", err)
	}

	return nil
}

func assembleChunks(tempDir string, fileName string) error {
	chunkFiles, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read chunk directory: %v", err)
	}

	// Sort the chunk files to ensure correct ordering
	sortChunkFiles(chunkFiles)

	// Create the final file to write the assembled chunks
	finalPath := filepath.Join(".", fileName)
	finalFile, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("failed to create final file: %v", err)
	}
	defer finalFile.Close()

	// Read each chunk file and write its content to the final file
	for _, chunkFile := range chunkFiles {
		chunkPath := filepath.Join(tempDir, chunkFile.Name())
		chunkData, err := ioutil.ReadFile(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to read chunk file: %v", err)
		}

		_, err = finalFile.Write(chunkData)
		if err != nil {
			return fmt.Errorf("failed to write chunk to final file: %v", err)
		}
	}

	// Remove the temporary directory and its contents
	err = os.RemoveAll(tempDir)
	if err != nil {
		return fmt.Errorf("failed to remove temp directory: %v", err)
	}

	return nil
}

func sortChunkFiles(chunkFiles []os.FileInfo) {
	// Implement sorting logic for chunk files
	// This can be based on their names or any other criteria
	// For simplicity, we assume that the chunk files are named with numerical order (e.g., chunk0, chunk1, chunk2, ...)
	// You may need to modify this logic based on your requirements
	sort.Slice(chunkFiles, func(i, j int) bool {
		return chunkFiles[i].Name() < chunkFiles[j].Name()
	})
}

func processChunks() {
	tempDir := "temp"
	chunkFiles, err := ioutil.ReadDir(tempDir)
	if err != nil {
		log.Fatal("Failed to read chunk directory:", err)
	}

	for _, chunkFile := range chunkFiles {
		chunkPath := filepath.Join(tempDir, chunkFile.Name())
		chunkData, err := ioutil.ReadFile(chunkPath)
		if err != nil {
			log.Println("Failed to read chunk file:", err)
			continue
		}

		fileName := strings.TrimSuffix(chunkFile.Name(), filepath.Ext(chunkFile.Name()))
		err = saveChunk(fileName, chunkData)
		if err != nil {
			log.Println("Failed to save chunk:", err)
		}
	}

	// Assemble the chunks into original files
	fileChunks, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal("Failed to read file chunks directory:", err)
	}

	for _, fileChunk := range fileChunks {
		fileName := fileChunk.Name()
		fileExt := filepath.Ext(fileName)
		if fileExt == ".dat" {
			err = assembleChunks(".", strings.TrimSuffix(fileName, fileExt))
			if err != nil {
				log.Println("Failed to assemble chunks for file:", fileName)
			} else {
				log.Println("File", fileName, "assembled successfully")
			}
		}
	}
}
