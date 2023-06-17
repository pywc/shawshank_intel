package util

import (
	"github.com/pywc/shawshank_intel/testers/http_tester"
	"github.com/pywc/shawshank_intel/testers/https_tester"
)

type Report struct {
	dns   dns_tester.DNSResult
	http  http_tester.HTTPResult
	https https_tester.HTTPSResult
	quic  quic_tester.QUICResult
}
