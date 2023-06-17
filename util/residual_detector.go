package util

import (
	"github.com/pywc/shawshank_intel/config"
	"time"
)

var AllResidualDetected []ResidualDetected

type ResidualDetected struct {
	component  string
	resultCode int
	duration   float64
}

func DetectResidual(ip string, port int, component string) (int, float64) {
	/*
		-1: init
		0: no residual
		1: destination ip
		2: destination ip + port
	*/
	// test 0
	conn, err := ConnectViaProxy(ip, port)
	if err == nil {
		conn.Close()
		return 0, 0.0
	}

	// test 1
	newPort := 80
	if port == 80 {
		newPort = 443
	}
	conn, err = ConnectViaProxy(ip, newPort)
	if err != nil {
		startTime := time.Now()
		endTime := time.Now()
		duration := 0.0

		for {
			conn, err := ConnectViaProxy(ip, newPort)
			endTime = time.Now()
			duration = float64(endTime.UnixMilli()-startTime.UnixMilli()) / 1000

			if err == nil {
				conn.Close()
				break
			} else if duration >= config.ResidualTestThreshold {
				break
			}
		}

		newDetected := ResidualDetected{
			component:  component,
			resultCode: 2,
			duration:   duration,
		}

		AllResidualDetected = append(AllResidualDetected, newDetected)
		return 1, config.ResidualTestThreshold
	} else {
		conn.Close()
	}

	// test 2
	conn, err = ConnectViaProxy(ip, port)
	if err != nil {
		startTime := time.Now()
		endTime := time.Now()
		duration := 0.0

		for {
			conn, err := ConnectViaProxy(ip, port)
			endTime = time.Now()
			duration = float64(endTime.UnixMilli()-startTime.UnixMilli()) / 1000

			if err == nil {
				conn.Close()
				break
			} else if duration >= config.ResidualTestThreshold {
				break
			}
		}

		newDetected := ResidualDetected{
			component:  component,
			resultCode: 2,
			duration:   duration,
		}

		AllResidualDetected = append(AllResidualDetected, newDetected)
		return 2, duration
	} else {
		conn.Close()
	}

	return -10, 0.0
}
