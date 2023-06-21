package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/testers/dns_tester"
	"github.com/pywc/shawshank_intel/testers/http_tester"
	"github.com/pywc/shawshank_intel/testers/https_tester"
	"github.com/pywc/shawshank_intel/testers/ip_tester"
	"github.com/pywc/shawshank_intel/testers/quic_tester"
	"github.com/pywc/shawshank_intel/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestReport struct {
	Country    string                   `json:"country"`
	ProxyIP    string                   `json:"proxy_ip"`
	HostDomain string                   `json:"host_domain"`
	HostIP     string                   `json:"host_ip"`
	Residual   []util.ResidualDetected  `json:"residual,omitempty"`
	DNS        dns_tester.DNSResult     `json:"dns"`
	IP         ip_tester.IPResult       `json:"ip"`
	HTTP       http_tester.HTTPResult   `json:"http"`
	HTTPS      https_tester.HTTPSResult `json:"https"`
	QUIC       quic_tester.QUICResult   `json:"quic"`
}

func TestDomain(country string, domain string, ip string) {
	util.PrintInfo(domain, "initiating tests...")
	util.ResetResidualDetector()

	report := TestReport{
		Country:    country,
		ProxyIP:    config.ProxyIP,
		HostDomain: domain,
		HostIP:     ip,
		DNS:        dns_tester.TestDNS(ip, domain),
		IP:         ip_tester.TestIP(ip),
		HTTP:       http_tester.TestHTTP(ip, domain),
		HTTPS:      https_tester.TestHTTPS(ip, domain),
		QUIC:       quic_tester.TestQUIC(ip, domain),
		Residual:   util.AllResidualDetected,
	}

	err := saveToDB(report)
	if err != nil {
		return
	}
}

func saveToDB(report TestReport) error {
	util.PrintInfo(report.HostDomain, "saving results...")
	j, _ := json.MarshalIndent(report, "", "    ")
	fmt.Println(string(j))

	client, err := connectDB()
	if err != nil {
		util.PrintError("", err)
		return err
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			util.PrintError("", err)
		}
	}(client, context.Background())

	collection := client.Database("shawshank").Collection("results")

	_, err = collection.InsertOne(context.Background(), report)
	if err != nil {
		util.PrintError("", err)
		return err
	}

	return nil
}

func connectDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		util.PrintError("", err)
		return nil, err
	}
	return client, nil
}
