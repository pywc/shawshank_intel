package throttle_tester

import (
	"github.com/influxdata/tdigest"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/testers/https_tester"
	"github.com/pywc/shawshank_intel/util"
	utls "github.com/refraction-networking/utls"
	"time"
)

var ConnectTimes = make([]float64, 0)

func AddConnectTime(t float64) {
	ConnectTimes = append(ConnectTimes, t)
}

func IsThrottled(t float64) bool {
	td := tdigest.NewWithCompression(config.TDigestCompression)
	for _, x := range ConnectTimes {
		td.Add(x, 1)
	}
	quantile := 1 - (config.ThrottlePValThreshold / 2)
	threshold := td.Quantile(quantile)

	if t > threshold {
		return true
	}

	return false
}

func InitThrottleChecker() error {
	testList := make([]string, 0)

	data, err := util.ReadCsvFile("./data/alexa_top_100.csv")
	if err != nil {
		return err
	}

	for _, row := range data {
		testList = append(testList, row[1])
	}

	for _, domain := range testList {
		// request configuration
		req := "GET / HTTP/1.1\r\n" +
			"Host: " + domain + "\r\n" +
			"Accept: */*\r\n" +
			"User-Agent: " + config.UserAgent + "\r\n" +
			"Connection: close\r\n\r\n"

		utlsConfig := utls.Config{
			ServerName: domain,
		}

		startTime := time.Now()
		_, _, err := https_tester.SendHTTPSRequest(domain, domain, 443, req, &utlsConfig)
		endTime := time.Now()
		if err != nil {
			continue
		}

		AddConnectTime(float64(endTime.UnixMilli()-startTime.UnixMilli()) / 1000)
	}

	return nil
}
