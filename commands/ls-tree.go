package commands

import (
	"encoding/hex"
	"fmt"
	"git/objects"
	"io"
	"os"
	"strings"
)

func LsTree(nameOnly *bool, treeHash string) (string, error) {
	_, _, bufReader := objects.Read(treeHash)

	for {
		header, err := bufReader.ReadBytes(0)
		if err != nil && err != io.EOF {
			fmt.Errorf("read header from .git/objects: %w", err)
			return "", err
		}

		if len(header) == 0 {
			break
		}

		header = header[:len(header)-1]

		modeName := string(header)

		buf := make([]byte, 20)

		limitedReader := &io.LimitedReader{R: bufReader, N: 20}
		n, err := io.ReadFull(limitedReader, buf)
		if err != nil {
			if err != io.ErrUnexpectedEOF {
				fmt.Println("Error reading bytes:", err, n)
				return "", err
			}
		}

		parts := strings.SplitN(modeName, " ", 2)

		mode := parts[0]
		name := parts[1]

		if *nameOnly == true {
			fmt.Println(name)

		} else {
			hashString := hex.EncodeToString(buf)
			kind, _, _ := objects.Read(hashString)
			fmt.Fprintf(os.Stdout, "%06s %s %s    %s", mode, kind, hashString, name)
		}

		fmt.Println("")
	}
	return "", nil
}
