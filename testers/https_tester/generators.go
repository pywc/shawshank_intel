package https_tester

import (
	"bytes"
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/jpillora/go-tld"
	"github.com/pywc/shawshank_intel/util"
	utls "github.com/refraction-networking/utls"
	"log"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type RequestWord struct {
	Servername   string
	CipherSuites []uint16
	MinVersion   uint16
	MaxVersion   uint16
	Certificate  []utls.Certificate
}

// Returns of an HTTP request for URL.
func CreateTLSConfig(requestWord RequestWord) *utls.Config {
	//Set max version to TLS 1.2 if we want cipher suites to be configurable. TLS 1.3 cipher suites are not configurable
	//https://pkg.go.dev/crypto/tls#Config
	if len(requestWord.CipherSuites) > 0 {
		if requestWord.MaxVersion == 772 || requestWord.MaxVersion == 0 {
			requestWord.MaxVersion = 771
		}
	}

	var cert utls.Certificate
	if len(requestWord.Certificate) > 0 {
		cert = requestWord.Certificate[0]
		c, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			log.Println("[https_fuzzer.CreateTLSConfig] Error parsing Certificate")
		}
		if c.Subject.CommonName == "" {
			clientCertPEM, ClientCertKey, err := GenerateCertificate(requestWord.Servername)
			if err != nil {
				log.Println("[https_fuzzer.CreateTLSConfig] Could not generate client Certificate")
				log.Println(err)
				return &utls.Config{
					ServerName:         requestWord.Servername,
					InsecureSkipVerify: true,
					CipherSuites:       requestWord.CipherSuites,
					MinVersion:         requestWord.MinVersion,
					MaxVersion:         requestWord.MaxVersion,
				}
			}
			cert, err = utls.X509KeyPair(clientCertPEM, ClientCertKey)
			if err != nil {
				log.Println("[https_fuzzer.CreateTLSConfig] Could not generate client Certificate")
				log.Println(err)
				return &utls.Config{
					ServerName:         requestWord.Servername,
					InsecureSkipVerify: true,
					CipherSuites:       requestWord.CipherSuites,
					MinVersion:         requestWord.MinVersion,
					MaxVersion:         requestWord.MaxVersion,
				}
			}
		}
	}

	//Handle servername changes - This has to be done at runtime, since the strategies would be selected first, but the servername itself is only known at runtime
	serverNameParts := strings.Split(requestWord.Servername, "|")
	if len(serverNameParts) > 1 {
		//ServerNameParts[1] contains the strategy to be run at runtime
		if serverNameParts[1] == "omit" {
			return &utls.Config{
				InsecureSkipVerify: true,
				CipherSuites:       requestWord.CipherSuites,
				MinVersion:         requestWord.MinVersion,
				MaxVersion:         requestWord.MaxVersion,
				Certificates:       []utls.Certificate{cert},
			}
		} else if serverNameParts[1] == "empty" {
			return &utls.Config{
				ServerName:         "",
				InsecureSkipVerify: true,
				CipherSuites:       requestWord.CipherSuites,
				MinVersion:         requestWord.MinVersion,
				MaxVersion:         requestWord.MaxVersion,
				Certificates:       []utls.Certificate{cert},
			}
		} else if serverNameParts[1] == "repeat" {
			//Now there should be a third part that says how many times to repeat
			repeatTimes, err := strconv.Atoi(serverNameParts[2])
			if err != nil {
				log.Println("[https_fuzzer.CreateTLSConfig] Error converting string into integer (repeat)")
				log.Println(err)
				log.Println("Reverting to default")
				return &utls.Config{
					ServerName:         requestWord.Servername,
					InsecureSkipVerify: true,
					CipherSuites:       requestWord.CipherSuites,
					MinVersion:         requestWord.MinVersion,
					MaxVersion:         requestWord.MaxVersion,
					Certificates:       []utls.Certificate{cert},
				}
			}
			serverName := util.Repeat(serverNameParts[0], repeatTimes)
			return &utls.Config{
				ServerName:         serverName,
				InsecureSkipVerify: true,
				CipherSuites:       requestWord.CipherSuites,
				MinVersion:         requestWord.MinVersion,
				MaxVersion:         requestWord.MaxVersion,
				Certificates:       []utls.Certificate{cert},
			}
		} else if serverNameParts[1] == "reverse" {
			return &utls.Config{
				ServerName:         util.Reverse(serverNameParts[0]),
				InsecureSkipVerify: true,
				CipherSuites:       requestWord.CipherSuites,
				MinVersion:         requestWord.MinVersion,
				MaxVersion:         requestWord.MaxVersion,
				Certificates:       []utls.Certificate{cert},
			}
		} else if serverNameParts[1] == "tld" {
			domainParts, _ := tld.Parse("https://" + serverNameParts[0])
			var serverName string
			if domainParts.Subdomain != "" {
				serverName = domainParts.Subdomain + "." + domainParts.Domain + "." + serverNameParts[2]
			} else {
				serverName = domainParts.Domain + "." + serverNameParts[2]
			}
			return &utls.Config{
				ServerName:         serverName,
				InsecureSkipVerify: true,
				CipherSuites:       requestWord.CipherSuites,
				MinVersion:         requestWord.MinVersion,
				MaxVersion:         requestWord.MaxVersion,
				Certificates:       []utls.Certificate{cert},
			}
		} else if serverNameParts[1] == "subdomain" {
			domainParts, _ := tld.Parse("https://" + serverNameParts[0])
			serverName := serverNameParts[2] + "." + domainParts.Domain + "." + domainParts.TLD
			return &utls.Config{
				ServerName:         serverName,
				InsecureSkipVerify: true,
				CipherSuites:       requestWord.CipherSuites,
				MinVersion:         requestWord.MinVersion,
				MaxVersion:         requestWord.MaxVersion,
				Certificates:       []utls.Certificate{cert},
			}
		}
	}

	return &utls.Config{
		ServerName:         requestWord.Servername,
		InsecureSkipVerify: true,
		CipherSuites:       requestWord.CipherSuites,
		MinVersion:         requestWord.MinVersion,
		MaxVersion:         requestWord.MaxVersion,
		Certificates:       []utls.Certificate{cert},
	}

}

func GenerateAlternatives(alternatives []string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	choice := rand.Intn(len(alternatives))
	return alternatives[choice]
}

func GenerateAllAlternatives(alternatives []string) []string {
	return alternatives
}

func GenerateHostNameRandomPadding() string {
	prefixPaddingLength := rand.Intn(5)
	suffixPaddingLength := rand.Intn(5)
	hostnameWithPadding := strings.Repeat("*", prefixPaddingLength)
	hostnameWithPadding += "%s"
	hostnameWithPadding += strings.Repeat("*", suffixPaddingLength)
	return hostnameWithPadding
}
func GenerateAllHostNamePaddings(hostname string) []string {
	var hostnameWithAllPadding []string
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			hostnameWithPadding := strings.Repeat("*", i)
			hostnameWithPadding += hostname
			hostnameWithPadding += strings.Repeat("*", j)
			hostnameWithAllPadding = append(hostnameWithAllPadding, hostnameWithPadding)
		}

	}
	return hostnameWithAllPadding
}

var versionAlternatives = []string{fmt.Sprint(tls.VersionTLS10), fmt.Sprint(tls.VersionTLS11), fmt.Sprint(tls.VersionTLS12), fmt.Sprint(tls.VersionTLS13)}

func GenerateVersionAlternatives() string {
	return GenerateAlternatives(versionAlternatives)
}

func GenerateAllVersionAlternatives() []string {
	return GenerateAllAlternatives(versionAlternatives)
}

var cipherSuiteAlternatives = []string{fmt.Sprint(tls.TLS_RSA_WITH_RC4_128_SHA),
	fmt.Sprint(tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA),
	fmt.Sprint(tls.TLS_RSA_WITH_AES_128_CBC_SHA),
	fmt.Sprint(tls.TLS_RSA_WITH_AES_256_CBC_SHA),
	fmt.Sprint(tls.TLS_RSA_WITH_AES_128_CBC_SHA256),
	fmt.Sprint(tls.TLS_RSA_WITH_AES_128_GCM_SHA256),
	fmt.Sprint(tls.TLS_RSA_WITH_AES_256_GCM_SHA384),
	fmt.Sprint(tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA),
	fmt.Sprint(tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA),
	fmt.Sprint(tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA),
	fmt.Sprint(tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256),
	fmt.Sprint(tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384),
	fmt.Sprint(tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384),
	fmt.Sprint(tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256),
	fmt.Sprint(tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256),
	fmt.Sprint(tls.TLS_AES_128_GCM_SHA256),
	fmt.Sprint(tls.TLS_AES_256_GCM_SHA384),
	fmt.Sprint(tls.TLS_CHACHA20_POLY1305_SHA256)}

func GenerateCipherSuiteAlternatives() string {
	return GenerateAlternatives(cipherSuiteAlternatives)
}

func GenerateAllCipherSuiteAlternatives() []string {
	return GenerateAllAlternatives(cipherSuiteAlternatives)
}

func GenerateCertificate(commonname string) ([]byte, []byte, error) {

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			Organization:  []string{"Internet Measurement"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test"},
			StreetAddress: []string{"Test 123"},
			PostalCode:    []string{"12345"},
			CommonName:    commonname,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand2.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	caBytes, err := x509.CreateCertificate(rand2.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Internet Measurement"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test"},
			StreetAddress: []string{"Test 123"},
			PostalCode:    []string{"12345"},
			CommonName:    commonname,
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand2.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand2.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return certPEM.Bytes(), certPrivKeyPEM.Bytes(), nil
}

var commonnames = []string{"www.test.com", "XXXX", "www.pornhub.com"}

func GenerateCertificateAlternatives() string {
	return GenerateAlternatives(commonnames)
}

func GenerateAllCertificateAlternatives() []string {
	return GenerateAllAlternatives(commonnames)
}

// https://azbigmedia.com/business/here-are-2021s-most-popular-tlds-and-domain-registration-trends/
var TLDs = []string{"%s|tld|com", "%s|tld|xyz", "%s|tld|net", "%s|tld|club", "%s|tld|me", "%s|tld|org", "%s|tld|co", "%s|tld|shop", "%s|tld|info", "%s|tld|live"}

// https://securitytrails.com/blog/most-popular-subdomains-mx-records#:~:text=As%20you%20can%20see%2C%20the,forums%2C%20wiki%2C%20community).
var Subdomains = []string{"%s|subdomain|www", "%s|subdomain|mail", "%s|subdomain|forum", "%s|subdomain|m", "%s|subdomain|blog", "%s|subdomain|shop", "%s|subdomain|forums", "%s|subdomain|wiki", "%s|subdomain|community", "%s|subdomain|ww1"}

var hostnames = []string{"%s|omit", "%s|empty", "%s|repeat|2", "%s|repeat|3", "%s|reverse"}

func GenerateTLDAlternatives() string {
	return GenerateAlternatives(TLDs)
}

func GenerateAllTLDAlternatives(hostname string) []string {
	u, _ := tld.Parse("https://" + hostname)
	tldAlternatives := GenerateAllAlternatives(TLDs)

	for i, alt := range tldAlternatives {
		tldAlternatives[i] = fmt.Sprintf(alt, u.Domain)
	}

	return tldAlternatives
}

func GenerateSubdomainsAlternatives() string {
	return GenerateAlternatives(Subdomains)
}

func GenerateAllSubdomainsAlternatives(hostname string) []string {
	subdomainAlternatives := GenerateAllAlternatives(Subdomains)

	for i, alt := range subdomainAlternatives {
		subdomainAlternatives[i] = fmt.Sprintf(alt, hostname)
	}

	return subdomainAlternatives
}

func GenerateHostNameAlternatives() string {
	return GenerateAlternatives(hostnames)
}

func GenerateAllHostNameAlternatives(hostname string) []string {
	hostnameAlternatives := GenerateAllAlternatives(hostnames)

	for i, alt := range hostnameAlternatives {
		hostnameAlternatives[i] = fmt.Sprintf(alt, hostname)
	}

	return hostnameAlternatives
}
