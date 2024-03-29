package main

import (
	"context"
	"fmt"
	"gotpb/gotpb"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Config file must be passed\n")
		return
	}

	conf_file := os.Args[1]
	log.Printf("Loading config from %s", conf_file)
	conf := gotpb.GetConf(conf_file)

	//gotpb.RunSingleCheck(conf)
	gotpb.InitRestService(conf)

	quit := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	log.Printf("Service running on: http://localhost:%d\nUse Ctrl^C to exit", conf.Port)

	go periodicCheck(quit, conf)

	server := http.Server{Addr: fmt.Sprintf(":%d", conf.Port)}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("%v", err)
		}
	}()
	<-sig
	quit <- true

	time.Sleep(time.Second)
	if err := server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

func periodicCheck(quit chan bool, conf gotpb.Config) {
	log.Printf("Periodic task started")

	for {
		select {
		case <-quit:
			log.Printf("Terminating peridic download")
			return
		case <-time.After(time.Duration(conf.Interval) * time.Hour):
			gotpb.RunSingleCheck(conf)
			log.Printf("Waiting for %d hours before next check", conf.Interval)
		}
	}
}
