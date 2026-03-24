package main

import (
	"io"
	"log"
	"os"

	"github.com/mistweaverco/kuba/cmd/kuba"
	"github.com/mistweaverco/kuba/internal/lib/fileutils"
)

func main() {
	f, err := os.OpenFile(fileutils.JoinPath(fileutils.GetTempPath(), "kuba.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("Warning: failed to close log file: %v", closeErr)
		}
	}()
	wrt := io.Writer(f)
	log.SetOutput(wrt)
	kuba.Execute()
}
