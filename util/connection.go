package util

import (
	"bufio"
	"crypto/x509"
	"fmt"
	"github.com/pywc/shawshank_intel/config"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
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
func SendHTTPTraffic(conn net.Conn, request string) (*http.Response, error) {
	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Failed to send HTTP request:", err)
		return nil, err
	}

	// Read the response
	response, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendHTTPSTraffic Send HTTP GET request with TLS and get response
func SendHTTPSTraffic(conn net.Conn, request string, utlsConfig *utls.Config,
	extensions []utls.TLSExtension, chloID utls.ClientHelloID) (*http.Response, error) {
	// Create a TLS connection over the proxy connection
	tlsConn := utls.UClient(conn, utlsConfig, chloID)
	tlsConn.BuildHandshakeState()

	// Add extensions
	for _, ext := range extensions {
		tlsConn.Extensions = append(tlsConn.Extensions, ext)
	}

	// Set a timeout for the TLS handshake
	tlsConn.SetDeadline(time.Now().Add(10 * time.Second))

	// Perform the TLS handshake
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}

	// Reset the deadline
	tlsConn.SetDeadline(time.Time{})

	_, err = tlsConn.Write([]byte(request))
	if err != nil {
		return nil, err
	}

	// Read the response
	response, err := http.ReadResponse(bufio.NewReader(tlsConn), nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type TLSSession struct {
	version      uint16
	cipherSuite  uint16
	ID           []byte
	masterSecret []byte
	serverCerts  []*x509.Certificate
}

// GetNewTLSSession
func GetNewTLSSession(conn net.Conn, request string, utlsConfig *utls.Config) (*TLSSession, error) {
	sess := TLSSession{}

	// Create a TLS connection over the proxy connection
	tlsConn := utls.UClient(conn, utlsConfig, utls.HelloGolang)

	err := tlsConn.BuildHandshakeState()
	if err != nil {
		return nil, err
	}

	tlsConn.HandshakeState.Hello.SessionId = nil

	// Set a timeout for the TLS handshake
	tlsConn.SetDeadline(time.Now().Add(10 * time.Second))

	// Perform the TLS handshake
	err = tlsConn.Handshake()
	if err != nil {
		return nil, err
	}

	// Reset the deadline
	tlsConn.SetDeadline(time.Time{})

	_, err = tlsConn.Write([]byte(request))
	if err != nil {
		return nil, err
	}

	state := tlsConn.HandshakeState

	sess.version = state.ServerHello.Vers
	sess.cipherSuite = state.ServerHello.CipherSuite
	sess.ID = state.ServerHello.SessionId
	sess.masterSecret = state.MasterSecret
	sess.serverCerts = tlsConn.ConnectionState().PeerCertificates

	return &sess, nil
}

// ResumeTLSSession
func ResumeTLSSession(conn net.Conn, request string, sess TLSSession) (*http.Response, error) {
	utlsConfig := utls.Config{
		ServerName:         config.DummyServerDomain,
		InsecureSkipVerify: true,
		MinVersion:         utls.VersionTLS12,
		MaxVersion:         utls.VersionTLS12,
	}

	// Create a TLS connection over the proxy connection
	tlsConn := utls.UClient(conn, &utlsConfig, utls.HelloGolang)

	// Set Session ID and State
	state := utls.MakeClientSessionState(
		nil,
		sess.version,
		sess.cipherSuite,
		sess.masterSecret,
		sess.serverCerts,
		nil,
	)

	tlsConn.HandshakeState.Hello.SessionId = sess.ID
	err := tlsConn.SetSessionState(state)
	if err != nil {
		return nil, err
	}

	// Set a timeout for the TLS handshake
	tlsConn.SetDeadline(time.Now().Add(10 * time.Second))

	// Perform the TLS handshake
	err = tlsConn.Handshake()
	if err != nil {
		return nil, err
	}

	// Reset the deadline
	tlsConn.SetDeadline(time.Time{})

	// Send HTTP request
	_, err = tlsConn.Write([]byte(request))
	if err != nil {
		return nil, err
	}

	// Read the response
	response, err := http.ReadResponse(bufio.NewReader(tlsConn), nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
