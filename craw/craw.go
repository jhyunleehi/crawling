package craw

import (
	"bufio"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//"github.com/PuerkitoBio/goquery"
	"github.com/anaskhan96/soup"
	log "github.com/sirupsen/logrus"
)

type Craw struct {
	host    string
	path    string
	url     string
	product string
	key     int
	rtitle  map[int]string
	rdate   map[int]string
	rnation map[int]string
	rstar   map[int]string
	rrating map[int]string
	rbody   map[int]string
}

func NewCraw(hostname, pathname, productname string) *Craw {
	trend := Craw{
		host:    hostname,
		path:    pathname,
		url:     hostname + pathname,
		product: productname,
		key:     1,
		rtitle:  map[int]string{},
		rdate:   map[int]string{},
		rnation: map[int]string{},
		rstar:   map[int]string{},
		rrating: map[int]string{},
		rbody:   map[int]string{},
	}
	return &trend
}

func (c *Craw) GetWebData(url string) error {
	//url := c.host + c.path
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

func (c *Craw) ParparseBody(body string) error {
	doc := soup.HTMLParse(string(body))
	div := doc.FindAll("div", "data-hook", "review")
	for _, d := range div {
		for _, a1 := range d.FindAll("a", "data-hook", "review-title") {
			a1 := a1.Find("span")
			title := a1.Text()
			log.Printf("title==> [%s]", title)
			c.rtitle[c.key] = title
		}
		a1 := d.Find("span", "data-hook", "review-date")
		t1 := a1.Text()
		t2 := strings.Split(t1, "on")
		t3 := strings.Trim(t2[1], " ")
		rdate, _ := time.Parse("January 02, 2006", t3)
		wdate := rdate.Format("2006-01-02")
		wnation := t2[0]
		c.rdate[c.key] = wdate
		c.rnation[c.key] = wnation
		log.Printf("dater==>[%s]", wdate)
		log.Printf("nation==> [%s]", wnation)
		for _, a1 := range d.FindAll("i", "data-hook", "review-star-rating") {
			a1 := a1.Find("span", "class", "a-icon-alt")
			star := string(a1.Text()[0])
			c.rstar[c.key] = star
			log.Printf("star==> [%s]", star)
		}
		for _, d1 := range d.FindAll("span", "data-hook", "review-body") {
			a2 := d1.Find("span")
			strbody := a2.Text()
			c.rbody[c.key] = strbody
			log.Printf("body==> [%s]", strbody)
		}
		c.key++
	}

	div = doc.FindAll("div", "id", "cm_cr-pagination_bar")

	for _, d := range div {
		alink := d.FindAll("li", "class", "a-last")
		for _, l := range alink {
			a1 := l.Find("a")
			if a1.NodeValue != "" {
				nextpage := a1.Attrs()["href"]
                newurl := c.host + nextpage
				c.url = newurl
				log.Debugf("[%s]", c.url)
				err := c.GetWebData(newurl)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	return nil
}

func (c *Craw) WriteToFile() error {
	filename := c.product + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	wr := csv.NewWriter(bufio.NewWriter(file))
	// csv 내용 쓰기
	for i := 1; i < c.key; i++ {
		wr.Write([]string{
			strconv.Itoa(i),
			c.rdate[i],
			c.rstar[i],
			c.rrating[i],
			c.rnation[i],
			c.rtitle[i],
			c.rbody[i],
		})
		wr.Flush()
	}
	return nil
}
