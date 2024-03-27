package objects

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
)

func Read(objectHash string) (string, int64, *bufio.Reader) {

	filePath := fmt.Sprintf(".git/objects/%s/%s", objectHash[:2], objectHash[2:])
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Errorf("open in .git/objects: %w", err)
	}
	zReader, err := zlib.NewReader(file)
	if err != nil {
		fmt.Errorf("create zlib reader: %w", err)
	}

  defer zReader.Close()

	// Wrap the zlib reader with a buffered reader
	bufReader := bufio.NewReader(zReader)

	// Read until the null terminator (0 byte)
	header, err := bufReader.ReadBytes(0)
	if err != nil && err != io.EOF {
		fmt.Errorf("read header from .git/objects: %w", err)
	}

	headerStr := string(header)

	parts := strings.SplitN(headerStr, " ", 2)

	if len(parts) != 2 {
		fmt.Errorf(".git/objects file header did not start with a known type: '%s'", headerStr)
	}

	kind := parts[0]
	size := 0
	sizeU, err := fmt.Sscanf(parts[1], "%d", &size)
	if err != nil {
		fmt.Errorf(".git/objects file header size is not valid: '%s'", parts[1])
	}

	return kind, int64(sizeU), bufReader

}
