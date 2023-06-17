package https_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	utls "github.com/refraction-networking/utls"
	"io"
	"strings"
)

type FilteredHTTPS struct {
	component  string
	resultCode int
}

// SendHTTPSRequest Returns result_code, response_body
/*
	Result Code Entry
	=========================
	-10: unhandled error
	-1: proxy error
	0: success
	1: reset
	2: refused
	3: silent drop
	4: TODO: throttle
	5: invalid certificate
*/
func SendHTTPSRequest(domain string, ip string, port int, req string, utlsConfig *utls.Config) (int, string, error) {
	// Fetch via proxy
	conn, err := util.ConnectViaProxy(ip, port)

	if err != nil {
		if strings.Contains(err.Error(), "general SOCKS server failure") {
			// cannot connect to proxy
			return -1, "", nil
		} else if ip == config.DummyServerIP {
			// residual censorship detection mode
			util.DetectResidual(ip, port, "https")
			conn, err = util.ConnectViaProxy(ip, port)
			if err != nil {
				return -10, "", err
			}
		} else {
			// unknown error
			return -10, "", err
		}
	}

	resp, err := util.SendHTTPSTraffic(conn, req, utlsConfig, nil, utls.HelloGolang)
	if err != nil {
		return -10, "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	conn.Close()

	// check tcp errors
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			// connection reset
			return 1, "", nil
		} else if strings.Contains(err.Error(), "connection refused") {
			// connection refused
			return 2, "", nil
		} else if strings.Contains(err.Error(), "i/o timeout") {
			// connection timeout
			return 3, "", nil
		} else if strings.Contains(err.Error(), "certificate is valid for") {
			return 5, "", nil
		} else {
			// unknown error
			return -10, "", err
		}
	}

	return 0, string(respBody), nil
}

// SendHTTPSRequestCustom Returns result_code, response_body
/*
	Result Code Entry
	=========================
	-10: unhandled error
	-1: proxy error
	0: success
	1: reset
	2: refused
	3: silent drop
	4: TODO: throttle
	5: invalid certificate
*/
func SendHTTPSRequestCustom(domain string, ip string, port int,
	req string, extensions []utls.TLSExtension) (int, string, error) {
	// Fetch via proxy
	conn, err := util.ConnectViaProxy(ip, port)

	if err != nil {
		if strings.Contains(err.Error(), "general SOCKS server failure") {
			// cannot connect to proxy
			return -1, "", nil
		} else if ip == config.ESNIIP {
			// residual censorship detection mode
			util.DetectResidual(ip, port, "https")
			conn, err = util.ConnectViaProxy(ip, port)
			if err != nil {
				return -10, "", err
			}
		} else {
			// unknown error
			return -10, "", err
		}
	}

	utlsConfig := utls.Config{
		InsecureSkipVerify: true,
		MinVersion:         utls.VersionTLS12,
		MaxVersion:         utls.VersionTLS13,
	}

	resp, err := util.SendHTTPSTraffic(conn, req, &utlsConfig, extensions, utls.HelloRandomizedNoALPN)
	if err != nil {
		return -10, "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	conn.Close()

	// check tcp errors
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			// connection reset
			return 1, "", nil
		} else if strings.Contains(err.Error(), "connection refused") {
			// connection refused
			return 2, "", nil
		} else if strings.Contains(err.Error(), "i/o timeout") {
			// connection timeout
			return 3, "", nil
		} else if strings.Contains(err.Error(), "certificate is valid for") {
			return 5, "", nil
		} else {
			// unknown error
			return -10, "", err
		}
	}

	return 0, string(respBody), nil
}
