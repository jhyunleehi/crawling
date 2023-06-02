package main

import (
	"crawling/craw"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
		NoColors:        true,
	})
}

func main() {
	var product string
	if len(os.Args) < 2 {
		//panic("에러: 2개 미만의 argument")
		product = "B084RGZ3P7"
	} else {
		product = os.Args[1]
	}
	//"B084RGZ3P7"
	host := "https://www.amazon.com"
	path := "/product-reviews/" + product
	c := craw.NewCraw(host, path, product)
	err := c.GetWebData()
	if err != nil {
		log.Error(err)
	}
	err = c.WriteToFile()
	if err != nil {
		log.Error(err)
	}
}
