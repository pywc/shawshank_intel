package https_tester

import (
	"github.com/pywc/shawshank_intel/testers/http_tester"
	"github.com/pywc/shawshank_intel/util"
	utls "github.com/refraction-networking/utls"
	"log"
	"strconv"
)

func FuzzSender(hostname string, ip string, req string, component string, utlsConfig *utls.Config) *FilteredHTTPS {
	resultCode, _, err := SendHTTPSRequest(hostname, ip, 443, req, utlsConfig)
	util.PrintInfo(hostname, component+" result: "+strconv.Itoa(resultCode))
	if resultCode == 0 {
		return nil
	} else if resultCode == -10 {
		log.Println("[*] Error - " + hostname + " - " + err.Error())
	}

	filtered := FilteredHTTPS{
		component:  component,
		resultCode: resultCode,
	}

	return &filtered
}

func CheckServerNamePadding(hostname string, ip string) []FilteredHTTPS {
	serverNameAllPadding := GenerateAllHostNamePaddings(hostname)
	filteredList := make([]FilteredHTTPS, len(serverNameAllPadding))

	for _, testComponent := range serverNameAllPadding {
		reqWord := RequestWord{
			Servername: testComponent,
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, testComponent, utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckMinVersionAlternate(hostname string, ip string) []FilteredHTTPS {
	minVerAllAlternate := GenerateAllVersionAlternatives()
	filteredList := make([]FilteredHTTPS, len(minVerAllAlternate))

	for _, testComponent := range minVerAllAlternate {
		minVersionInt, _ := strconv.ParseUint(testComponent, 10, 16)
		reqWord := RequestWord{
			Servername: hostname,
			MinVersion: uint16(minVersionInt),
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, testComponent, utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckMaxVersionAlternate(hostname string, ip string) []FilteredHTTPS {
	maxVerAllAlternate := GenerateAllVersionAlternatives()
	filteredList := make([]FilteredHTTPS, len(maxVerAllAlternate))

	for _, testComponent := range maxVerAllAlternate {
		maxVersionInt, _ := strconv.ParseUint(testComponent, 10, 16)
		reqWord := RequestWord{
			Servername: hostname,
			MaxVersion: uint16(maxVersionInt),
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, testComponent, utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckCipherSuiteAlternate(hostname string, ip string) []FilteredHTTPS {
	cipherSuiteAllAlternate := GenerateCipherSuiteAlternatives()
	filteredList := make([]FilteredHTTPS, len(cipherSuiteAllAlternate))

	for _, testComponent := range cipherSuiteAllAlternate {
		reqWord := RequestWord{
			Servername:   hostname,
			CipherSuites: []uint16{uint16(testComponent)},
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, string(testComponent), utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckClientCertAlternate(hostname string, ip string) []FilteredHTTPS {
	clientCertAllAlternate := GenerateAllCertificateAlternatives()
	filteredList := make([]FilteredHTTPS, len(clientCertAllAlternate))

	for _, testComponent := range clientCertAllAlternate {
		clientCertPEM, ClientCertKey, _ := GenerateCertificate(testComponent)
		clientCert, _ := utls.X509KeyPair(clientCertPEM, ClientCertKey)

		reqWord := RequestWord{
			Servername:  hostname,
			Certificate: []utls.Certificate{clientCert},
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, testComponent, utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckServernameAlternate(hostname string, ip string) []FilteredHTTPS {
	servernameAllAlternate := GenerateAllHostNameAlternatives(hostname)
	filteredList := make([]FilteredHTTPS, len(servernameAllAlternate))

	for _, testComponent := range servernameAllAlternate {
		reqWord := RequestWord{
			Servername: testComponent,
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, testComponent, utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckServernameTLDAlternate(hostname string, ip string) []FilteredHTTPS {
	servernameTLDAllAlternate := GenerateAllTLDAlternatives(hostname)
	filteredList := make([]FilteredHTTPS, len(servernameTLDAllAlternate))

	for _, testComponent := range servernameTLDAllAlternate {
		reqWord := RequestWord{
			Servername: testComponent,
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, testComponent, utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckServernameSubdomainAlternate(hostname string, ip string) []FilteredHTTPS {
	servernameAllAlternate := GenerateAllSubdomainsAlternatives(hostname)
	filteredList := make([]FilteredHTTPS, len(servernameAllAlternate))

	for _, testComponent := range servernameAllAlternate {
		reqWord := RequestWord{
			Servername: testComponent,
		}

		utlsConfig := CreateTLSConfig(reqWord)

		req := http_tester.FormatHttpRequest(http_tester.RequestWord{Hostname: hostname})
		filtered := FuzzSender(hostname, ip, req, testComponent, utlsConfig)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}
