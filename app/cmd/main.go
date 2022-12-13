package main

import (
	"GoAsyncJofogasParcer/internal/app"
	"fmt"
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
	t := time.Now()
	app.Run()

	fmt.Println(time.Since(t))
}
