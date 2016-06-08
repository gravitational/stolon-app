package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	sentinels := flag.Int("sentinels", 2, "number of sentinels")
	proxies := flag.Int("proxies", 2, "number of proxies")
	password := flag.String("password", "password1", "initial database user password")

	flag.Parse()

	log.Infof("starting stolonboot")
	err := bootCluster(*sentinels, *proxies, *password)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
