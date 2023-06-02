package main

import (
	"crawling/craw"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

type TSuite struct {
	suite.Suite
	mycraw *craw.Craw
}


func TestSuite(t *testing.T) {
	suite.Run(t, new(TSuite))
}

func (s *TSuite) Test_review() {	
	product := "B0B6HRB7T4"	
	//https://www.amazon.com/product-reviews/B084RGZ3P7	
	host := "https://www.amazon.com"
	path := "/product-reviews/" + product	
	url := host + path
	s.mycraw = craw.NewCraw(host, path, product)
	err := s.mycraw.GetWebData(url)
	if err != nil {
		log.Error(err)
	}
	err = s.mycraw.WriteToFile()
	if err != nil {
		log.Error(err)
	}
}

