//usr/bin/env go run "$0" "$@"; exit

package main

import (
	"github.com/pywc/shawshank_intel/db"
	"github.com/pywc/shawshank_intel/util"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: ./main.go <PROXY_CSV>")
		return
	}

	log.Println("Welcome to Shawshank Intel!")

	proxyList, err := util.FetchProxy(os.Args[1])
	if err != nil {
		return
	}

	// test for each proxy
	for _, proxy := range proxyList {
		// get parameters
		countryCode := proxy[0]
		proxyType := proxy[1]
		proxyIP := proxy[2]
		proxyPort, _ := strconv.Atoi(proxy[3])
		proxyUsername := proxy[4]
		proxyPassword := proxy[5]

		// set proxy
		util.SetProxy(proxyIP, proxyPort, proxyUsername, proxyPassword, proxyType)

		// get domains
		testList, err := util.GetTestDomains(countryCode)
		if err != nil {
			continue
		}

		// test for each domain
		for _, domain := range testList {
			ip, err := util.ResolveIPLocally(domain)
			if err != nil {
				continue
			}

			db.TestDomain(countryCode, domain, ip)
		}
	}
}
