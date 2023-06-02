package craw

import (
	"fmt"
	"io"
	"net/http"

	//"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

func (c *Craw) GetColly(url string) error {
	c.co = colly.NewCollector()

	// Find and visit all links
	c.co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.co.Visit("http://go-colly.org/")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Thunder Client (https://www.thunderclient.com)")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return nil
	}
	//fmt.Println(string(body))
	err = c.ParparseBody(string(body))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
