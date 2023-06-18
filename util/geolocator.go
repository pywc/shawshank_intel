package util

import (
	"log"
	"net/url"
)

// TODO: implement maxmind geolocator
func GetCountry(ip string) string {
	return "us"
}

func GetTestDomains(countryCode string) ([]string, error) {
	log.Println("Fetching test domains...")

	testList := make([]string, 0)

	data, err := ReadCsvFile("./test-lists/lists/" + countryCode + ".csv")

	urlList := make([]string, 0)
	for _, row := range data {
		urlEntry, _ := url.Parse(row[0])
		urlList = append(urlList, urlEntry.Host)
	}

	urlList = RemoveDuplicateStr(urlList)

	log.Println("Fetched " + string(len(urlList)) + " domains")

	return testList, err
}
