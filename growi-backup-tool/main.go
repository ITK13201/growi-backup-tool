package main

import (
	"github.com/itk13201/growi-backup-tool/cmd"
	"log"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
