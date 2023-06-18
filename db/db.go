package db

import (
	"context"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/testers/dns_tester"
	"github.com/pywc/shawshank_intel/testers/http_tester"
	"github.com/pywc/shawshank_intel/testers/https_tester"
	"github.com/pywc/shawshank_intel/testers/quic_tester"
	"github.com/pywc/shawshank_intel/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestReport struct {
	country    string
	proxyIP    string
	hostDomain string
	hostIP     string
	dns        dns_tester.DNSResult
	http       http_tester.HTTPResult
	https      https_tester.HTTPSResult
	quic       quic_tester.QUICResult
}

func TestDomain(country string, domain string) {
	util.PrintInfo(domain, "start testing...")

	hostIP, err := util.ResolveIPLocally(domain)
	if err != nil {
		return
	}

	report := TestReport{
		country:    country,
		proxyIP:    config.ProxyIP,
		hostDomain: domain,
		hostIP:     hostIP,
		dns:        dns_tester.TestDNS("", domain),
		http:       http_tester.TestHTTP("", domain),
		https:      https_tester.TestHTTPS("", domain),
		quic:       quic_tester.TestQUIC("", domain),
	}

	err = saveToDB(report)
	if err != nil {
		return
	}
}

func saveToDB(report TestReport) error {
	util.PrintInfo(report.hostDomain, "saving results...")
	client, err := connectDB()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("shawshank").Collection("results")
	_, err = collection.InsertOne(context.Background(), report)
	if err != nil {
		return err
	}

	return nil
}

func connectDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		util.PrintError(config.ProxyIP, "", err)
		return nil, err
	}
	return client, nil
}
