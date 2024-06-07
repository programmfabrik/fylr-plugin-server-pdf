package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
)

const chrome string = "chromium"

type proxyWriter struct {
	data   []byte
	proxy  io.Writer
	portCh chan int
}

func newProxyWriter(w io.Writer) *proxyWriter {
	pw := proxyWriter{
		proxy:  w,
		portCh: make(chan int),
	}
	return &pw
}

func (pw *proxyWriter) Write(data []byte) (n int, err error) {
	n, err = pw.proxy.Write(data)
	if err == nil && pw.portCh != nil {
		pw.data = append(pw.data, data...)
		s := string(pw.data)
		re := regexp.MustCompile(`DevTools listening on ws://0\.0\.0\.0:([0-9]+)`)
		matches := re.FindAllStringSubmatch(s, -1)
		if len(matches) > 0 {
			port, _ := strconv.Atoi(matches[0][1])
			pw.portCh <- port
			close(pw.portCh)
			pw.portCh = nil
		}
	}
	return n, err
}

func startChrome(ctx context.Context, prog string) (port int, process *os.Process, exit chan bool, err error) {
	_, err = exec.LookPath(prog)
	if err != nil {
		return 0, nil, nil, err
	}
	cmd := exec.CommandContext(ctx, prog,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--remote-debugging-address=0.0.0.0",
		"--remote-debugging-port=0",
	)
	wg := sync.WaitGroup{}
	pw := newProxyWriter(os.Stderr)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case port = <-pw.portCh:
			log.Printf("%s running on port %d", prog, port)
		}
	}()
	cmd.Stderr = pw
	cmd.Stdout = os.Stdout
	exit = make(chan bool, 1)
	go func() {
		defer func() {
			exit <- true
		}()
		err := cmd.Start()
		if err != nil {
			log.Fatalf("%s start failed with error: %s", prog, err.Error())
		}
		process = cmd.Process
		log.Printf("started %s, waiting to finish...", prog)
		err = cmd.Wait()
		if err != nil {
			log.Printf("%s exited with error: %s", prog, err.Error())
			return
		}
		log.Printf("%s exited", prog)
	}()
	wg.Wait()
	return port, process, exit, nil
}
