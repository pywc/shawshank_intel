package util

import (
	"errors"
	"github.com/pywc/shawshank_intel/config"
	"log"
	"net/url"
	"strconv"
)

// TODO: implement maxmind geolocator
func GetCountry(ip string) string {
	return "us"
}

func GetTestDomains(countryCode string) ([]string, error) {
	log.Println("Fetching test domains...")

	data, err := ReadCsvFile("./test-lists/lists/" + countryCode + ".csv")
	if err != nil {
		PrintError(config.ProxyIP, "", errors.New("failed to load test list"))
		return nil, err
	}
	data = data[1:]

	urlList := make([]string, 0)
	for _, row := range data {
		urlEntry, _ := url.Parse(row[0])
		urlList = append(urlList, urlEntry.Host)
	}

	urlList = RemoveDuplicateStr(urlList)

	log.Println("Fetched " + strconv.Itoa(len(urlList)) + " domains")

	return urlList, nil
}
