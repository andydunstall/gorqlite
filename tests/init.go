package tests

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logging enabled")
	}

	rand.Seed(time.Now().Unix())
}
