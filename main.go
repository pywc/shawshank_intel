package main

import (
	"fmt"
	"github.com/pywc/shawshank_intel/testers/http_tester"
)

func main() {
	//result := -1
	//util.SetProxy("116.209.68.219", 12345, "", "")
	//result = ip_tester.TestPort("142.250.190.36", 80)
	resultCode, redirectURL := http_tester.CheckHTTPConnectivity("naver.com", "223.130.195.200")
	fmt.Println(resultCode, redirectURL)
}
