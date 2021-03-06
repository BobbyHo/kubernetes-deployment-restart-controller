package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/xing/kubernetes-deployment-restart-controller/src/controller"
	"github.com/xing/kubernetes-deployment-restart-controller/src/util"
)

var options struct {
	RestartCheckPeriod int  `short:"c" long:"restart-check-period" env:"RESTART_CHECK_PERIOD" description:"Time interval to check for pending restarts in milliseconds" default:"500"`
	RestartGracePeriod int  `short:"r" long:"restart-grace-period" env:"RESTART_GRACE_PERIOD" description:"Time interval to compact restarts in seconds" default:"5"`
	Verbose            int  `short:"v" long:"verbose" env:"VERBOSE" description:"Be verbose"`
	Version            bool `long:"version" description:"Print version information and exit"`
}

// VERSION represents the current version of the release.
const VERSION = "v1.1.0"

func main() {
	util.ParseArgs(&options)

	if options.Version {
		printVersion()
		return
	}

	http.Handle("/metrics", prometheus.Handler())
	addr := fmt.Sprintf("0.0.0.0:10254")
	go func() { glog.Fatal(http.ListenAndServe(addr, nil)) }()

	controller := controller.NewDeploymentConfigController(time.Duration(options.RestartCheckPeriod)*time.Millisecond, time.Duration(options.RestartGracePeriod)*time.Second)
	util.InstallSignalHandler(controller.Stop)

	err := controller.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Controller terminated: %s", err)
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("kubernetes-deployment-restart-controller %s %s/%s %s\n", VERSION, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
