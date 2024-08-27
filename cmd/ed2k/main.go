package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"official/md4"
	"os"
	"path/filepath"
)

const chunkSize = 9500 * 1024

func calculateChunkHashes(filePath string) ([][]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunkHashes [][]byte
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		hasher := md4.New()
		hasher.Write(buffer[:n])
		chunkHashes = append(chunkHashes, hasher.Sum(nil))

		if err == io.EOF {
			break
		}
	}

	return chunkHashes, nil
}

func computeEd2kHash(filePath string) (string, error) {
	chunkHashes, err := calculateChunkHashes(filePath)
	if err != nil {
		return "", err
	}

	if len(chunkHashes) == 1 {
		// Single chunk (file size <= 9500KB)
		return hex.EncodeToString(chunkHashes[0]), nil
	}

	// Multiple chunks, concatenate the chunk hashes and hash them
	finalHasher := md4.New()
	for _, chunkHash := range chunkHashes {
		finalHasher.Write(chunkHash)
	}
	finalHash := finalHasher.Sum(nil)
	return hex.EncodeToString(finalHash), nil
}

func generateEd2kLink(filePath string) (string, error) {
	fileName := filepath.Base(filePath)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	ed2kHash, err := computeEd2kHash(filePath)
	if err != nil {
		return "", err
	}

	ed2kLink := fmt.Sprintf("ed2k://|file|%s|%d|%s|/", fileName, fileSize, ed2kHash)
	return ed2kLink, nil
}

func main() {
	filePath := os.Args[1]

	ed2kLink, err := generateEd2kLink(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("eD2k Address:", ed2kLink)
}
