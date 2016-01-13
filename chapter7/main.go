package main

import (
	"github.com/qiniu/log"
)

func main() {

	conf := Config{
		Addr:    ":7090",
		MgoAddr: "127.0.0.1:27017",
		MgoDB:   "short",
		MgoColl: "short",
	}

	log.Fatal(Run(conf))
}
