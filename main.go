package main

import (
	"fmt"
	"github.com/pywc/shawshank_intel/testers/https_tester"
	"github.com/pywc/shawshank_intel/util"
)

func main() {
	util.SetProxy("116.209.68.219", 12345, "", "")
	//result := ip_tester.TestPort("142.250.190.36", 80)
	result := https_tester.CheckHTTPSConnectivity("pornhub.com", "66.254.114.41")
	fmt.Println(result)
}
