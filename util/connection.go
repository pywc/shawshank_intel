package util

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// Parse the Content-Length header from the response headers
func getContentLength(headers string) int64 {
	lines := strings.Split(headers, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Content-Length:") {
			value := strings.TrimSpace(line[len("Content-Length:"):])
			length, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				return length
			}
		}
	}
	return 0
}

// Check if the response is chunked encoded
func isChunkedEncoding(headers string) bool {
	lines := strings.Split(headers, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Transfer-Encoding:") {
			value := strings.TrimSpace(line[len("Transfer-Encoding:"):])
			return strings.ToLower(value) == "chunked"
		}
	}
	return false
}

// Read the response body when it is chunked encoded
func readChunkedBody(reader *bufio.Reader) (string, error) {
	body := ""
	for {
		// Read the chunk size
		sizeLine, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		sizeLine = strings.TrimSpace(sizeLine)
		chunkSize, err := strconv.ParseInt(sizeLine, 16, 64)
		if err != nil {
			return "", err
		}

		if chunkSize == 0 {
			// Reached the end of chunked response
			break
		}

		// Read the chunk data
		chunk := make([]byte, chunkSize)
		_, err = io.ReadFull(reader, chunk)
		if err != nil {
			return "", err
		}

		// Append the chunk data to the body
		body += string(chunk)

		// Read the CRLF after the chunk data
		_, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
	}

	return body, nil
}

// Read the response body when it has a known length
func readResponseBody(reader *bufio.Reader, contentLength int64) (string, error) {
	body := make([]byte, contentLength)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ConnectNormally Connect to the address normally
func ConnectNormally(addr string, port int) (net.Conn, error) {
	// Connect to the SOCKS5 proxy
	conn, err := net.Dial("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		fmt.Println("Failed to normally connect:", err)
		return nil, err
	}

	return conn, nil
}

// ConnectViaProxy Connect to the address through the proxy
func ConnectViaProxy(addr string, port int) (net.Conn, error) {
	// Create a SOCKS5 proxy dialer
	dialer, err := proxy.SOCKS5("tcp", ParseProxy(), nil, proxy.Direct)
	if err != nil {
		fmt.Println("Failed to create proxy dialer:", err)
		return nil, err
	}

	// Connect to the SOCKS5 proxy
	proxyConn, err := dialer.Dial("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		fmt.Println("Failed to connect to the SOCKS5 proxy:", err)
		return nil, err
	}

	return proxyConn, nil
}

// SendHTTPTraffic Send HTTP GET request and get response
func SendHTTPTraffic(conn net.Conn, request string) (string, error) {
	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Failed to send HTTP request:", err)
		return "", err
	}

	// Read the response header
	respHeader := ""
	reader := bufio.NewReader(conn)
	for {
		packet, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		respHeader += packet

		// Check if the response headers are complete
		if strings.TrimSpace(packet) == "" {
			break
		}
	}

	// Read the response body
	isChunked := isChunkedEncoding(respHeader)
	respBody := ""
	if isChunked {
		respBody, err = readChunkedBody(reader)
		if err != nil {
			return "", err
		}
	} else {
		// Read the response body with Content-Length header
		contentLength := getContentLength(respHeader)
		respBody, err = readResponseBody(reader, contentLength)
		if err != nil {
			return "", err
		}
	}

	return respHeader + respBody, nil
}

// SendHTTPSTraffic Send HTTP GET request with TLS and get response
func SendHTTPSTraffic(conn net.Conn, request string, tlsConfig *tls.Config) (string, error) {
	// TODO: return session ticket

	// Create a TLS connection over the proxy connection
	tlsConn := tls.Client(conn, tlsConfig)

	// Set a timeout for the TLS handshake
	tlsConn.SetDeadline(time.Now().Add(10 * time.Second))

	// Perform the TLS handshake
	err := tlsConn.Handshake()
	if err != nil {
		return "", err
	}

	// Reset the deadline
	tlsConn.SetDeadline(time.Time{})

	_, err = tlsConn.Write([]byte(request))
	if err != nil {
		return "", err
	}

	// Read the response header
	respHeader := ""
	reader := bufio.NewReader(tlsConn)
	for {
		packet, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		respHeader += packet

		// Check if the response headers are complete
		if strings.TrimSpace(packet) == "" {
			break
		}
	}

	// Read the response body
	isChunked := isChunkedEncoding(respHeader)
	respBody := ""
	if isChunked {
		respBody, err = readChunkedBody(reader)
		if err != nil {
			return "", err
		}
	} else {
		// Read the response body with Content-Length header
		contentLength := getContentLength(respHeader)
		respBody, err = readResponseBody(reader, contentLength)
		if err != nil {
			return "", err
		}
	}

	return respHeader + respBody, nil
}
