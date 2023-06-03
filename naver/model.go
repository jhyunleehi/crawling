package naver

type Request struct {
	Page            int    `json:"page"`            //: 5,
	PageSize        int    `json:"pageSize"`        //: 20,
	MerchantNo      float64 `json:"merchantNo"`      //: "500054903",
	OriginProductNo string `json:"originProductNo"` //: "3806830917",
	SortType        string `json:"sortType"`        //: "REVIEW_RANKING"
}

type Contents struct {
	Id                   string `json:"id"`                   //: "4265112063",
	Type                 string `json:"reviewType"`           //: "NORMAL",
	ServiceType          string `json:"reviewServiceType"`    //: "SELLBLOG",
	ReviewScore          int    `json:"reviewScore"`          //": 5,
	ReviewContent        string `json:"reviewContent"`        //: "살짝 익으니 더맛있네요",
	CreateData           string `json:"createDate"`           //: "2023-06-03T01:30:23.966+00:00",
	ProductNo            string `json:"productNo"`            //: "3813375295",
	ProductName          string `json:"productName"`          //: "말바우시장 선김치 전라도 국산 알타리 총각 김치 2kg",
	ProductOptionContent string `json:"productOptionContent"` //: "착한 총각김치: 2kg",
}

type Review struct {
	ContentsList  []Contents `json:"contents"`       //
	Page          int        `json:"page"`          //: 5,
	Size          int        `json:"size"`          //: 20,
	TotalElements int        `json:"totalElemens"` //: 33168,
	Totalpages    int        `json:"totalPages"`    //: 1659,
	First         bool       `json:"first"`         //: false,
	Last          bool       `json:"last"`          //: false
}

type Crawl struct {
	rcount  int
	rtitle  map[int]string
	rdate   map[int]string
	rnation map[int]string
	rstar   map[int]string
	rrating map[int]string
	rbody   map[int]string
}

func NewCrawl() *Crawl {
	c := Crawl{
		rtitle:  map[int]string{},
		rdate:   map[int]string{},
		rnation: map[int]string{},
		rstar:   map[int]string{},
		rrating: map[int]string{},
		rbody:   map[int]string{},
	}
	return &c
}

/*
{
	"name": "말바우시장 선김치 전라도 국산 알타리 총각 김치 2kg",
	"@context": "https://schema.org",
	"@type": "Product",
	"image": "https://shop-phinf.pstatic.net/20230517_113/168428312740492GnN_JPEG/1325607795912398_862297780.jpg",
	"description": "말바우시장 박인영",
	"sku": 500054903,
	"mpn": "3813375295",
	"productID": "3813375295",
	"category": "식품>김치>총각김치",
	"offers": {
	  "@type": "Offer",
	  "price": 19000,
	  "priceCurrency": "KRW",
	  "availability": "http://schema.org/InStock",
	  "url": "https://shopping.naver.com/outlink/itemdetail/3813375295"
	},
	"brand": {
	  "@type": "Brand",
	  "logo": "http://shop1.phinf.naver.net/20181102_169/CM10357_1541152062472gaL3k_JPEG/64459315027094704_1150752429.jpg",
	  "slogan": "착한김치 선김치"
	},
	"aggregateRating": {
	  "@type": "AggregateRating",
	  "bestRating": 5,
	  "worstRating": 1,
	  "ratingValue": 4.72,
	  "reviewCount": 33168,
	  "ratingCount": 33168
	}
  }
*/
