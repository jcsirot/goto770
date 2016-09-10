package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/golang/glog"

	"github.com/jcsirot/goto770/core"
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
	glog.Infoln("Starting GoTo7/70")
	defer glog.Flush()
	wg.Add(1)
	go func() {
		defer wg.Done()
		core.Start()
	}()
	wg.Wait()
	runtime.LockOSThread()
}
