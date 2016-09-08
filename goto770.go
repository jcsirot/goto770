package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
	"runtime"

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
	glog.Infoln("Starting GoTo7/70")
	defer glog.Flush()
	go core.Start()
	runtime.LockOSThread()
}
