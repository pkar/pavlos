package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"api"
	"github.com/golang/glog"
	"models"
)

func main() {
	env := flag.String("env", "development", "environment")
	port := flag.String("port", "8000", "port")
	dbPath := flag.String("dbPath", "localhost", "path to mongo db")
	loadArticles := flag.Bool("load", false, "path to db")
	maxProcs := flag.Int("maxProcs", runtime.NumCPU(), "number of process")
	cpuProfile := flag.Bool("cpuProfile", false, "enable cpu profiling")
	memProfile := flag.Bool("memProfile", false, "enable memory profiling")
	flag.Parse()

	runtime.GOMAXPROCS(*maxProcs)

	if *cpuProfile {
		f, err := os.Create("mbp.prof")
		if err == nil {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	if *memProfile {
		f, err := os.Create("mbp.mprof")
		if err == nil {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}

	glog.Info(*env, *port, *dbPath)

	dbInstance, err := models.New(*dbPath)
	if err != nil {
		glog.Fatal(err)
	}
	if *loadArticles {
		dbInstance.InitCategories()
		dbInstance.ItemUpdater()
		os.Exit(0)
	}

	go dbInstance.ItemUpdaterCron()

	apiServer, err := api.New(*port, *env, dbInstance)
	if err != nil {
		glog.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case sig := <-interrupt:
		apiServer.Close()
		glog.Infof("Captured %v, exiting...", sig)
	}
}
