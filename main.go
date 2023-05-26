package main

import (
	"fmt"
	"github.com/pywc/shawshank_intel/testers/ip_tester"
)

func main() {
	result := ip_tester.TestPort("142.250.190.36", 80)
	fmt.Println(result)
}
