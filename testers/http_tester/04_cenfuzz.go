package http_tester

import (
	"log"
)

func FuzzSender(hostname string, ip string, req string, component string) *FilteredHTTP {
	resultCode, _, redirectURL, err := SendHTTPRequest(hostname, ip, 80, req)
	if resultCode == 0 {
		return nil
	} else if resultCode == -10 {
		log.Println("[*] Error - " + hostname + " - " + err.Error())
	}

	filtered := FilteredHTTP{
		component:   component,
		resultCode:  resultCode,
		redirectURL: redirectURL,
	}

	return &filtered
}

func CheckHostnamePadding(hostname string, ip string) []FilteredHTTP {
	hostnameAllPadding := GenerateAllHostNamePaddings(hostname)
	filteredList := make([]FilteredHTTP, len(hostnameAllPadding))

	for _, testComponent := range hostnameAllPadding {
		reqWord := RequestWord{
			Hostname: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckGetWordCapitalize(hostname string, ip string) []FilteredHTTP {
	getWords := GenerateAllCapitalizedPermutations("GET")
	filteredList := make([]FilteredHTTP, len(getWords))

	for _, testComponent := range getWords {
		reqWord := RequestWord{
			Hostname: hostname,
			GetWord:  testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckGetWordRemove(hostname string, ip string) []FilteredHTTP {
	getWordAllRemove := GenerateAllSubstringPermutations("GET")
	filteredList := make([]FilteredHTTP, len(getWordAllRemove))

	for _, testComponent := range getWordAllRemove {
		reqWord := RequestWord{
			Hostname: hostname,
			GetWord:  testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckGetWordAlternate(hostname string, ip string) []FilteredHTTP {
	getWordAllAlternate := GenerateAllGetAlternatives()
	filteredList := make([]FilteredHTTP, len(getWordAllAlternate))

	for _, testComponent := range getWordAllAlternate {
		reqWord := RequestWord{
			Hostname: hostname,
			GetWord:  testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHTTPWordCapitalize(hostname string, ip string) []FilteredHTTP {
	httpWords := GenerateAllCapitalizedPermutations("HTTP/1.1")
	filteredList := make([]FilteredHTTP, len(httpWords))

	for _, testComponent := range httpWords {
		reqWord := RequestWord{
			Hostname: hostname,
			HttpWord: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHTTPWordRemove(hostname string, ip string) []FilteredHTTP {
	httpWordAllRemove := GenerateAllSubstringPermutations("HTTP/1.1")
	filteredList := make([]FilteredHTTP, len(httpWordAllRemove))

	for _, testComponent := range httpWordAllRemove {
		reqWord := RequestWord{
			Hostname: hostname,
			HttpWord: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHTTPWordAlternate(hostname string, ip string) []FilteredHTTP {
	httpWordAllAlternate := GenerateAllHttpAlternatives()
	filteredList := make([]FilteredHTTP, len(httpWordAllAlternate))

	for _, testComponent := range httpWordAllAlternate {
		reqWord := RequestWord{
			Hostname: hostname,
			HttpWord: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHostWordCapitalize(hostname string, ip string) []FilteredHTTP {
	hostWords := GenerateAllCapitalizedPermutations("Host:")
	filteredList := make([]FilteredHTTP, len(hostWords))

	for _, testComponent := range hostWords {
		reqWord := RequestWord{
			Hostname: hostname,
			HostWord: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHostWordRemove(hostname string, ip string) []FilteredHTTP {
	hostWordAllRemove := GenerateAllSubstringPermutations("Host:")
	filteredList := make([]FilteredHTTP, len(hostWordAllRemove))

	for _, testComponent := range hostWordAllRemove {
		reqWord := RequestWord{
			Hostname: hostname,
			HostWord: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHostWordAlternate(hostname string, ip string) []FilteredHTTP {
	hostWordAllAlternate := GenerateAllHostAlternatives()
	filteredList := make([]FilteredHTTP, len(hostWordAllAlternate))

	for _, testComponent := range hostWordAllAlternate {
		reqWord := RequestWord{
			Hostname: hostname,
			HostWord: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHTTPDelimiterWordRemove(hostname string, ip string) []FilteredHTTP {
	httpDelimiterWordAllRemove := GenerateAllSubstringPermutations("\r\n")
	filteredList := make([]FilteredHTTP, len(httpDelimiterWordAllRemove))

	for _, testComponent := range httpDelimiterWordAllRemove {
		reqWord := RequestWord{
			Hostname:          hostname,
			HttpDelimiterWord: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckPathAlternate(hostname string, ip string) []FilteredHTTP {
	pathAllAlternate := GenerateAllPathAlternatives()
	filteredList := make([]FilteredHTTP, len(pathAllAlternate))

	for _, testComponent := range pathAllAlternate {
		reqWord := RequestWord{
			Hostname: hostname,
			Path:     testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHeaderAlternate(hostname string, ip string) []FilteredHTTP {
	headerAllAlternate := GenerateAllHeaderAlternatives()
	filteredList := make([]FilteredHTTP, len(headerAllAlternate))

	for _, testComponent := range headerAllAlternate {
		reqWord := RequestWord{
			Hostname: hostname,
			Header:   testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHostnameAlternate(hostname string, ip string) []FilteredHTTP {
	hostnameAllAlternate := GenerateAllHostNameAlternatives(hostname)
	filteredList := make([]FilteredHTTP, len(hostnameAllAlternate))

	for _, testComponent := range hostnameAllAlternate {
		reqWord := RequestWord{
			Hostname: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHostnameTLDAlternate(hostname string, ip string) []FilteredHTTP {
	hostnameTLDAllAlternate := GenerateAllTLDAlternatives(hostname)
	filteredList := make([]FilteredHTTP, len(hostnameTLDAllAlternate))

	for _, testComponent := range hostnameTLDAllAlternate {
		reqWord := RequestWord{
			Hostname: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}

func CheckHostnameSubdomainAlternate(hostname string, ip string) []FilteredHTTP {
	subdomainAllAlternate := GenerateAllHostNameAlternatives(hostname)
	filteredList := make([]FilteredHTTP, len(subdomainAllAlternate))

	for _, testComponent := range subdomainAllAlternate {
		reqWord := RequestWord{
			Hostname: testComponent,
		}

		req := FormatHttpRequest(reqWord)
		filtered := FuzzSender(hostname, ip, req, testComponent)
		if filtered == nil {
			continue
		}

		filteredList = append(filteredList, *filtered)
	}

	return filteredList
}
