package util

import (
	"context"
	"errors"
	"net"
	"time"
)

func ResolveIPLocally(domain string) (string, error) {
	PrintInfo(domain, "resolving IP from Cloudflare...")

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}
	ipList, err := r.LookupHost(context.Background(), domain)

	if err != nil || len(ipList) == 0 {
		newErr := errors.New("cannot resolve IP from Cloudflare")
		PrintError(domain, newErr)
		return "", newErr
	}

	ip := ipList[0]

	PrintInfo(domain, "resolved IP "+ip)
	return ip, nil
}
