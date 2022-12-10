package main

import (
	"gotpb/gotpb"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Config file must be passed\n")
		return
	}

	conf_file := os.Args[1]
	log.Printf("Loading config from %s", conf_file)
	conf := gotpb.GetConf(conf_file)
	gotpb.RunSingleCheck(conf)
}
