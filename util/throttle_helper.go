package util

import (
	"github.com/influxdata/tdigest"
	"github.com/pywc/shawshank_intel/config"
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
