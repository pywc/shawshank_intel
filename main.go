package main

import (
	"./testers/ip"
	"fmt"
)

func main() {
	result := ip.TestPort("142.250.190.36", 80)
	fmt.Println(result)
}
