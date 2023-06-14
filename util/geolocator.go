package util

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
)

// TODO: implement maxmind geolocator
func GetCountry(ip string) string {
	return "us"
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, err
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// TODO: implement getting test domains
func GetTestDomains(countryCode string) ([]string, error) {
	testList := make([]string, 0)

	data, err := readCsvFile("./test-lists/lists/" + countryCode + ".csv")

	urlList := make([]string, 0)
	for _, row := range data {
		urlEntry, _ := url.Parse(row[0])
		urlList = append(urlList, urlEntry.Host)
	}

	urlList = removeDuplicateStr(urlList)

	for _, u := range urlList {
		fmt.Println(u)
	}

	return testList, err
}
