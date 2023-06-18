package util

import "log"

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
