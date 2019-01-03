package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/jcsirot/goto770/pkg/core"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func main() {
	var wg sync.WaitGroup
	log.Infoln("Starting GoTo7/70")
	wg.Add(1)
	go func() {
		defer wg.Done()
		core.Start()
	}()
	wg.Wait()
	runtime.LockOSThread()
}
