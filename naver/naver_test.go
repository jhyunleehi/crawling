package naver

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

type TSuite struct {
	suite.Suite
	r *Resty
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TSuite))
}

func (s *TSuite) Test_review() {

	//url := "https://shopping.naver.com/window-products/localfood/3813375295"
	url :="https://shopping.naver.com/logistics/products/8571752135"
	product := ""
	//https://shopping.naver.com/v1/reviews/paged-reviews
	host := "https://shopping.naver.com"
	path := "/v1/reviews/paged-reviews"
	//url := host + path
	s.r = NewResty(host, path, product, url)
	err := s.r.GetBaseData()
	err = s.r.GetReview()
	if err != nil {
		log.Error(err)
		return
	}
	if err != nil {
		log.Error(err)
	}
	// err = s.r.WriteToFile()
	// if err != nil {
	// 	log.Error(err)
	// }
}

//"https://shopping.naver.com/v1/reviews/paged-reviews/3806830917"
