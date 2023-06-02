package craw


import (


    "bufio"


    "encoding/csv"


    "io"


    "net/http"


    "os"


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


func (c *Craw) GetWebData() error {


    url := c.host + c.path


    req, _ := http.NewRequest("GET", url, nil)


    req.Header.Add("Accept", "*/*")


    req.Header.Add("User-Agent", "Thunder Client (https://www.thunderclient.com)")


    res, _ := http.DefaultClient.Do(req)


    defer res.Body.Close()


    body, _ := io.ReadAll(res.Body)


    //fmt.Println(string(body))


    err := c.ParparseBody(string(body))


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


        alink := d.Find("li", "class", "a-last")


        a1 := alink.Find("a")


        if a1.NodeValue != "" {


            newpath := a1.Attrs()["href"]


            c.path = newpath


            log.Debugf("[%s]", c.path)


            err := c.GetWebData()


            if err != nil {


                log.Error(err)


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


        wr.Write([]string{"A", "0.25"})


        wr.Write([]string{"B", "55.70"})


        wr.Flush()


    }


    return nil


}

