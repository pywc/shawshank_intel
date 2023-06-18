package util

import (
	"encoding/csv"
	"github.com/pywc/shawshank_intel/config"
	"log"
	"os"
)

func ReadCsvFile(filePath string) ([][]string, error) {
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

func RemoveDuplicateStr(strSlice []string) []string {
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

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Repeat(s string, n int) string {
	returnString := ""
	for i := 0; i < n; i++ {
		returnString += s
	}
	return returnString
}

func PrintError(proxyIP string, domain string, err error) {
	log.Println(proxyIP + " - " + err.Error() + " - " + domain)
}

func PrintInfo(domain string, info string) {
	log.Println(config.ProxyIP + " - " + domain + " - " + info)
}
