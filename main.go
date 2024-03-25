package main

import (
	"flag"
	"fmt"
	"git/commands"
	"os"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	catFileCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)
	catFilePrettyPrint := catFileCmd.Bool("p", false, "enable")

	writeObjectCmd := flag.NewFlagSet("hash-object", flag.ExitOnError)
	writeObject := writeObjectCmd.Bool("w", false, "enable")

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'init'")
		// Create necessary directories and files for initialization
		createDir(".git")
		createDir(".git/objects")
		createDir(".git/refs")
		createFile(".git/HEAD", "ref: refs/heads/main\n")
		fmt.Println("Initialized git directory")

	case "cat-file":
		catFileCmd.Parse(os.Args[2:])
		commands.CatFile(catFilePrettyPrint, catFileCmd.Args()[0])
	case "hash-object":
		// Implement hash-object command
    writeObjectCmd.Parse(os.Args[2:])
		commands.HashObject(writeObject, writeObjectCmd.Args()[0])
	}
}

func createDir(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Printf("Error creating directory %s: %s\n", path, err)
		os.Exit(1)
	}
}

func createFile(filename, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", filename, err)
		os.Exit(1)
	}
}
