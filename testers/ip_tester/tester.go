package ip_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
)

type IPResult struct {
	Handshake80  int `json:"handshake_80"`
	Handshake443 int `json:"handshake_443"`
}

func TestIP(ip string) IPResult {
	config.CurrentComponent = "ip"
	util.PrintInfo(ip, "testing IP...")

	result := IPResult{
		Handshake80:  CheckTCPHandshake(ip, 80),
		Handshake443: CheckTCPHandshake(ip, 443),
	}

	return result
}
