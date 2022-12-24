package main

import (
	"GoAsyncJofogasParcer/internal/app"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
}

func main() {
	for {
		timer1 := time.NewTimer(2 * time.Second)
		<-timer1.C
		app.Run()
	}
}
