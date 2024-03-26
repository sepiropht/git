package commands

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func writeBlob(file string, writer io.Writer) (string, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", file, err)
	}

	zlibWriter := zlib.NewWriter(writer)
	hasher := sha1.New()

	if _, err := fmt.Fprintf(hasher, "blob %d\x00", stat.Size()); err != nil {
		return "", err
	}

	fileReader, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", file, err)
	}

	defer fileReader.Close()

	multiWriter := io.MultiWriter(zlibWriter, hasher)

	if _, err := io.Copy(multiWriter, fileReader); err != nil {
		return "", fmt.Errorf("stream file into blob: %w", err)
	}

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}

func HashObject(writeFile *bool, file string) (string, error) {

	tmpFile, err := os.CreateTemp("", "tmp")
	if err != nil {
		return "", fmt.Errorf("Error creating temporary file: %v\n", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	hash, err := writeBlob(file, tmpFile)
	if err != nil {
		return "", fmt.Errorf("Error writing blob object: %v\n", err)
	}

	if *writeFile == true {
		hashDir := fmt.Sprintf(".git/objects/%s/", hash[:2])
		if err := os.MkdirAll(hashDir, 0755); err != nil {
			return "", fmt.Errorf("Error creating subdir of .git/objects: %v\n", err)
		}

		destination := filepath.Join(hashDir, hash[2:])
		if err := os.Rename(tmpFile.Name(), destination); err != nil {
			return "", fmt.Errorf("Error moving blob file into .git/objects: %v\n", err)
		}
	}
	fmt.Println(hash)
	return "", nil
}
