package main

import (
	"financeSys/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.WithFields(log.Fields{"location": "main", "method": "main"}).Fatal(err)
	}
}