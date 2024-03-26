package commands

import (
	"fmt"
	"git/objects"
	"io"
	"os"
)

func CatFile(prettyPrint *bool, objectHash string) (string, error) {
	/* fmt.Println("cat-file executed")
	   fmt.Println("  prettyPrint:", *prettyPrint)
	   fmt.Println("  objectHash:", objectHash)
	*/
	if *prettyPrint == false {
		return "", fmt.Errorf("mode must be given without -p, and we don't support mode")
	}

	kind, size, bufReader := objects.Read(objectHash)

	limitedReader := &io.LimitedReader{R: bufReader, N: size}

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
