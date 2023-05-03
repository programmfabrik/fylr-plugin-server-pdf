package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/programmfabrik/golib"
)

func endWithError(err error) {
	log.Print(err.Error())
	sendEvent(event{Type: "SERVER_PDF_GENERATE_ERROR", Info: map[string]any{
		"error": err.Error(),
	}})
	os.Exit(1)
}

var info infoT

func main() {
	prog := golib.GetEnv("")["SERVER_PDF_CHROME"]
	if prog == "" {
		prog = "chromium"
	}

	infoS := flag.String("info", "", "JSON with callback info sent by fylr")
	flag.Parse()

	if infoS != nil && *infoS != "" {
		err := json.Unmarshal([]byte(*infoS), &info)
		if err != nil {
			endWithError(err)
		}
	}

	var err error

	defer func() {
		log.Printf("done")
	}()

	body := pdfCreatorBody{}
	ctx := context.Background()
	dec := json.NewDecoder(os.Stdin)
	dec.DisallowUnknownFields()

	err = dec.Decode(&body)
	if err != nil {
		endWithError(err)
	}

	timeStartup := time.Now()

	port, process, exit, err := startChrome(ctx, prog)
	if err != nil {
		endWithError(err)
	}

	tookStartup := time.Since(timeStartup)

	f, err := os.CreateTemp("", "*.html")
	if err != nil {
		endWithError(err)
	}

	timeProduce := time.Now()

	_, err = f.Write([]byte(body.Document))
	if err != nil {
		endWithError(err)
	}

	log.Printf("opened file %s", f.Name())

	log.Printf("display header footer %t", *body.Properties.DisplayHeaderFooter)

	data, err := createPdf(ctx, "file://"+f.Name(), port, body.Properties)
	if err != nil {
		endWithError(err)
		return
	}

	tookProduce := time.Since(timeProduce)

	n, err := os.Stdout.Write(data)
	if err != nil {
		endWithError(err)
		return
	}
	sendEvent(event{Type: "SERVER_PDF_GENERATE", Info: map[string]any{
		"time startup":     tookStartup.String(),
		"time pdf produce": tookProduce.String(),
		"pdf size":         golib.HumanByteSize(uint64(n)),
	}})
	log.Printf("%d bytes written to stdout", n)
	process.Signal(syscall.SIGTERM)
	// wait for chrome to exit
	<-exit

	// remove temp html file
	err = os.Remove(f.Name())
	if err != nil {
		endWithError(err)
	}
	log.Printf("removed file %s", f.Name())

}
