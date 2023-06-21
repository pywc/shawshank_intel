package util

import (
	"encoding/json"
	"errors"
	"github.com/pywc/shawshank_intel/config"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Response struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	ASN     struct {
		ASNum   int    `json:"asnum"`
		OrgName string `json:"org_name"`
	} `json:"asn"`
	Geo struct {
		City       string  `json:"city"`
		Region     string  `json:"region"`
		RegionName string  `json:"region_name"`
		PostalCode string  `json:"postal_code"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
		TimeZone   string  `json:"tz"`
		LumCity    string  `json:"lum_city"`
		LumRegion  string  `json:"lum_region"`
	} `json:"geo"`
}

func FetchISP() error {
	conn, err := ConnectViaProxy("lumtest.com", 80, "http")

	req := "GET http://lumtest.com/myip.json HTTP/1.1\r\n" +
		"Host: lumtest.com\r\n" +
		"Accept: */*\r\n" +
		"User-Agent: " + config.UserAgent + "\r\n"
	if config.ProxyUsername != "" {
		req += "Proxy-Authorization: Basic " + ParseAuth() + "\r\n"
	}
	req += "\r\n"

	resp, err := SendHTTPTraffic(conn, req)
	var response Response

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	conn.Close()

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		PrintError("", err)
		return err
	}

	config.ProxyISP = response.ASN.OrgName
	config.ProxyCountry = strings.ToLower(response.Country)

	return nil
}

func GetTestDomains(countryCode string) ([]string, error) {
	log.Println("Fetching test domains...")

	data, err := ReadCsvFile("./test-lists/lists/" + countryCode + ".csv")
	if err != nil {
		PrintError("", errors.New("failed to load test list"))
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
