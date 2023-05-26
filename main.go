package main

import (
	"fmt"
	"github.com/pywc/shawshank_intel/testers/ip"
)

func main() {
	result := ip.TestPort("142.250.190.36", 80)
	fmt.Println(result)
}
