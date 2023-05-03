package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/programmfabrik/golib"
)

func main() {
	prog := golib.GetEnv("")["SERVER_PDF_CHROME"]
	if prog == "" {
		prog = "chromium"
	}

	body := pdfCreatorBody{}
	ctx := context.Background()
	dec := json.NewDecoder(os.Stdin)
	dec.DisallowUnknownFields()
	err := dec.Decode(&body)
	if err != nil {
		log.Fatal(err)
	}

	port, process, err := startChrome(ctx, prog)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.CreateTemp("", "*.html")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = os.Remove(f.Name())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("removed file %s", f.Name())
	}()

	_, err = f.Write([]byte(body.Document))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("opened file %s", f.Name())

	data, err := createPdf(ctx, "file://"+f.Name(), port, body.Properties)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	n, err := os.Stdout.Write(data)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Printf("%d bytes written to stdout", n)
	process.Signal(syscall.SIGTERM)
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
	// chrome exits the program
}
