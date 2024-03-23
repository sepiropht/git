package commands

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
)

func CatFile(prettyPrint *bool, objectHash string) (string, error) {
	/* fmt.Println("cat-file executed")
	   fmt.Println("  prettyPrint:", *prettyPrint)
	   fmt.Println("  objectHash:", objectHash)
	*/
	if *prettyPrint == false {
		return "", fmt.Errorf("mode must be given without -p, and we don't support mode")
	}

	filePath := fmt.Sprintf(".git/objects/%s/%s", objectHash[:2], objectHash[2:])
	file, err := os.Open(filePath)

	if err != nil {
		return "", fmt.Errorf("open in .git/objects: %w", err)
	}

	zReader, err := zlib.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("create zlib reader: %w", err)
	}

	defer zReader.Close()

	// Wrap the zlib reader with a buffered reader
	bufReader := bufio.NewReader(zReader)

	// Read until the null terminator (0 byte)
	header, err := bufReader.ReadBytes(0)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read header from .git/objects: %w", err)
	}

	headerStr := string(header)

	parts := strings.SplitN(headerStr, " ", 2)

	if len(parts) != 2 {
		return "", fmt.Errorf(".git/objects file header did not start with a known type: '%s'", headerStr)
	}

	kind := parts[0]
	size := 0
	sizeU, err := fmt.Sscanf(parts[1], "%d", &size)
	if err != nil {
		return "", fmt.Errorf(".git/objects file header size is not valid: '%s'", parts[1])
	}

	limitedReader := &io.LimitedReader{R: bufReader, N: int64(sizeU)}

	switch kind {
	case "blob":
		// Write the decompressed content to stdout
		_, err := io.Copy(os.Stdout, limitedReader.R)

		if err != nil {
			fmt.Println("Error writing to stdout:", err)
		}
	default:
		fmt.Printf("We do not yet know how to handle '%s'\n", kind)
	}

	return "", nil
}
