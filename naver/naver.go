package naver

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"

	nested "github.com/antonfisher/nested-logrus-formatter"
	resty "github.com/go-resty/resty/v2"
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

type Resty struct {
	name            string
	host            string
	path            string
	url             string
	sku             float64
	originProductNo string
	product         string
	pagekey         int
	client          *resty.Client
	request         *resty.Request
	crawl           Crawl
}

func NewResty(host, path, product, url string) *Resty {
	r := Resty{
		name:    product,
		host:    host,
		path:    path,
		url:     url,
		pagekey: 1,
	}
	r.client = resty.New()
	r.crawl = *NewCrawl()
	//r.client.SetDebug(true)
	r.client.SetDebug(false)
	r.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	r.client.SetTimeout(1 * time.Minute)
	r.request = r.client.R()
	return &r
}

func (r *Resty) GetBaseData() error {
	//first_url := "https://shopping.naver.com/window-products/localfood/" + r.product
	first_url := r.url
	soup.Header("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	soup.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit 537.36 (KHTML, like Gecko) Chrome")
	resp, err := soup.Get(first_url)
	if err != nil {
		log.Error(err)
		return err
	}
	//log.Debugf("%+v", resp)
	doc := soup.HTMLParse(resp)
	script_text := doc.FindAll("script")

	// for k, _ := range script_text {
	// 	script := script_text[k].Text()
	// 	log.Debugf("%+v", script)
	// }

	basescript := script_text[0].Text()
	log.Debugf("%s", basescript)
	var jres map[string]interface{}
	json.Unmarshal([]byte(basescript), &jres)
	if jres == nil {
		return nil
	}
	if val, ok := jres["sku"]; ok {
		r.sku = val.(float64)
	} else {
		log.Error("Not Found sku key")
		return errors.New("not found sku key ")
	}

	for _, item := range doc.FindAll("a") {
		log.Debug(item)
		c := item.Attrs()["class"]
		log.Debug(c)
		if strings.Contains(c, "a:rfd.report") {
			url := item.Attrs()["href"]
			token := strings.Split(url, "productNo=")
			if len(token) >= 2 {
				r.originProductNo = token[1]
				log.Debug(r.originProductNo)
				return nil
			}
		}
	}

	if r.originProductNo == "" {
		log.Error("not found origin product id")
		return errors.New("not found origin product id")
	}

	return nil
}

//<div class="-QJD7ahi1z"><span class="_2vERiTjT4T">신고센터</span><p class="_1hT0ZFItpQ">네이버㈜는 소비자의 보호와 사이트의 안전거래를 위해<br>신고센터를 운영하고 있습니다. 안전거래를 저해하는 경우 신고하여 주시기 바랍니다.</p>
//<a href="https://help.pay.naver.com/mail/form.nhn?alias=checkout_accuse&amp;productNo=8530881994"
// target="_blank" class="_3HTXTIRBSp N=a:rfd.report" role="button">신고하기</a></div>

func (r *Resty) GetReview() error {
	err := r.GetReviewAll()
	if err != nil {
		log.Error(err)
		return err
	}

	err = r.WriteToFile()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (r *Resty) GetReviewAll() error {
	reiewurl := "https://shopping.naver.com/v1/reviews/paged-reviews"
	log.Debugf("[%f][%s]", r.sku, r.product)
	in := Request{
		Page:            r.pagekey,
		PageSize:        20,
		MerchantNo:      r.sku,             // "500054903",
		OriginProductNo: r.originProductNo, //
		SortType:        "REVIEW_RANKING",
	}
	out := Review{}
	err := r.CallRestApi("POST", reiewurl, &in, &out)
	if err != nil {
		log.Error(err)
		return err
	}
	for _, contents := range out.ContentsList {
		r.crawl.rdate[r.crawl.rcount] = contents.CreateData
		r.crawl.rbody[r.crawl.rcount] = contents.ReviewContent
		r.crawl.rstar[r.crawl.rcount] = strconv.Itoa(contents.ReviewScore)
		r.crawl.rcount++
	}
	if out.Last {
		return nil
	} else {
		r.GetReviewAll()
	}

	return nil
}

func (r *Resty) CallRestApi(method, url string, in, out interface{}) (err error) {
	//log.Debugf("[%+v]", restitem)
	r.request = r.request.SetHeader("Accept", "*/*")
	r.request = r.request.SetHeader("User-Agent", "Thunder Client (https://www.thunderclient.com)")
	//url = "https://shopping.naver.com/v1/reviews/paged-reviews"
	if in != nil {
		body, _ := json.Marshal(in)
		r.request.SetBody(body)
	}

	var res *resty.Response
	switch method {
	case "GET":
		res, err = r.request.Get(url)
	case "POST":
		res, err = r.request.Post(url)
	case "PUT":
		res, err = r.request.Put(url)
	case "PATCH":
		res, err = r.request.Patch(url)
	case "DELETE":
		res, err = r.request.Delete(url)
	default:
		return errors.New("error")
	}
	if err != nil {
		log.Error(err)
		return err
	}
	if len(res.Body()) > 0 {
		log.Debugf("%s", string(res.Body()))
		json.Unmarshal(res.Body(), out)
	}
	return nil
}

func (r *Resty) WriteToFile() error {
	filename := r.product + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	wr := csv.NewWriter(bufio.NewWriter(file))
	// csv 내용 쓰기
	for i := 1; i < r.crawl.rcount; i++ {
		wr.Write([]string{
			strconv.Itoa(i),
			r.crawl.rdate[i],
			r.crawl.rstar[i],
			r.crawl.rrating[i],
			r.crawl.rnation[i],
			r.crawl.rtitle[i],
			r.crawl.rbody[i],
		})
		wr.Flush()
	}
	return nil
}
