//usr/bin/env go run "$0" "$@"; exit

package main

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/db"
	"github.com/pywc/shawshank_intel/util"
	"log"
	"math/rand"
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
		log.Println("testing country " + countryCode + "...")

		// set proxy
		util.SetProxy(proxyIP, proxyPort, proxyUsername, proxyPassword, proxyType)
		err := util.FetchISP()
		if err != nil {
			continue
		}

		// get domains
		testList, err := util.GetTestDomains(countryCode)
		if err != nil {
			continue
		}

		// Fisher–Yates shuffle
		for i := len(testList) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			testList[i], testList[j] = testList[j], testList[i]
		}

		if len(testList) < config.TestCount {
			log.Println("testing only " + strconv.Itoa(len(testList)) + " domains")
			testList = testList[:len(testList)]
		} else {
			log.Println("testing only " + strconv.Itoa(config.TestCount) + " domains")
			testList = testList[:30]
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
