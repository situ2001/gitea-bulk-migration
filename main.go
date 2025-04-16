package main

import (
	"github/situ2001.com/gitea-bulk-migration/cmd"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err := os.OpenFile("migration.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	log.Println("Starting gitea-bulk-migration...")

	cmd.Execute()

	log.Println("gitea-bulk-migration finished.")
}
