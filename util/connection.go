package util

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"strconv"
	"time"
)

func ConnectViaProxy(addr string, port int) (net.Conn, error) {
	// Create a SOCKS5 proxy dialer
	dialer, err := proxy.SOCKS5("tcp", ParseProxy(), nil, proxy.Direct)
	if err != nil {
		log.Fatal("Failed to create proxy dialer:", err)
		return nil, err
	}

	// Connect to the SOCKS5 proxy
	proxyConn, err := dialer.Dial("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		log.Fatal("Failed to connect to the SOCKS5 proxy:", err)
		return nil, err
	}

	return proxyConn, nil
}

func SendHTTPTraffic(conn net.Conn, request string) (string, error) {
	_, err := conn.Write([]byte(request))
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
		return "", err
	}

	// Example: Read the HTTP response
	buffer := make([]byte, 100000)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal("Failed to read HTTP response:", err)
		return "", err
	}

	return string(buffer[:n]), nil
}

func SendHTTPSTraffic(conn net.Conn, request string, tlsConfig tls.Config) (string, error) {
	// Create a TLS connection over the proxy connection
	tlsConn := tls.Client(conn, &tlsConfig)

	// Set a timeout for the TLS handshake
	tlsConn.SetDeadline(time.Now().Add(10 * time.Second))

	// Perform the TLS handshake
	err := tlsConn.Handshake()
	if err != nil {
		log.Fatal("TLS handshake failed:", err)
	}

	// Reset the deadline
	tlsConn.SetDeadline(time.Time{})

	_, err = tlsConn.Write([]byte(request))
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
		return "", err
	}

	// Example: Read the HTTP response
	buffer := make([]byte, 100000)
	n, err := tlsConn.Read(buffer)
	if err != nil {
		log.Fatal("Failed to read HTTP response:", err)
		return "", err
	}

	return string(buffer[:n]), nil
}
