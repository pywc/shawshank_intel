package main

import (
	"crypto/tls"
	"log"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	// Set up the SOCKS5 proxy address
	proxyAddress := "35.188.143.33:2408"

	// Set up the destination address and port
	destinationAddress := "www.naver.com"
	destinationPort := "443"

	// Create a SOCKS5 proxy dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		log.Fatal("Failed to create proxy dialer:", err)
	}

	// Connect to the SOCKS5 proxy
	proxyConn, err := dialer.Dial("tcp", net.JoinHostPort(destinationAddress, destinationPort))
	if err != nil {
		log.Fatal("Failed to connect to the SOCKS5 proxy:", err)
	}
	defer proxyConn.Close()

	// Create a TLS connection over the proxy connection
	tlsConn := tls.Client(proxyConn, &tls.Config{
		InsecureSkipVerify: true, // Skip certificate verification for simplicity
		ServerName:         destinationAddress,
	})

	// Set a timeout for the TLS handshake
	tlsConn.SetDeadline(time.Now().Add(10 * time.Second))

	// Perform the TLS handshake
	err = tlsConn.Handshake()
	if err != nil {
		log.Fatal("TLS handshake failed:", err)
	}

	// Reset the deadline
	tlsConn.SetDeadline(time.Time{})

	// Now you can use the tlsConn to send and receive encrypted data

	// Example: Send an HTTP GET request
	_, err = tlsConn.Write([]byte("GET / HTTP/1.1\r\nHost: " + destinationAddress + "\r\n\r\n"))
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}

	// Example: Read the HTTP response
	buffer := make([]byte, 100000)
	n, err := tlsConn.Read(buffer)
	if err != nil {
		log.Fatal("Failed to read HTTP response:", err)
	}

	// Print the HTTP response
	log.Println(string(buffer[:n]))
}