package main

import (
	"github.com/influxdata/tdigest"
	"log"
)

func main() {
	td := tdigest.NewWithCompression(1000)
	for _, x := range []float64{1.1, 2.2, 3.3, 4.4, 5.5, 5.5, 4.4, 3.3, 22, 1.1} {
		td.Add(x, 1)
	}

	// Compute Quantiles
	log.Println(td[0])
	log.Println("99th", td.Quantile(0.975))
}
