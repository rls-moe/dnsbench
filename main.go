package main

import (
	"fmt"

	"github.com/schollz/progressbar"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	numRequests      = kingpin.Flag("i", "Number of Requests per DNS server").Default("1000").Short('i').Uint16()
	wildcardEndpoint = kingpin.Flag("e", "DNS Wildcard Endpoint to benchmark against").Default("google.com").Short('e').String()
	dnsServers       = kingpin.Arg("dns servers", "DNS Servers to ping").Default("127.0.0.1", "8.8.8.8").Strings()
	sleepTimeout     = kingpin.Flag("wait-request", "Time to wait between requests in milliseconds").Default("100").Uint()
	antiCache        = kingpin.Flag("anti-cache", "Prepend randomized subdomains to query to prevent some caching. THIS REQUIRES A WILDCARD DNS ENTRY!").Default("false").Bool()
)

func main() {
	kingpin.Parse()
	for k := range *dnsServers {
		fmt.Printf("Testing %s...\n", (*dnsServers)[k])
		bar := progressbar.New(int(*numRequests))
		res, err := bench((*dnsServers)[k], *wildcardEndpoint, *numRequests, func(i, _ uint16) {
			bar.Set(int(i))
		})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Printf("\n"+
			"\tP00.5  = % 6.3fms\n"+
			"\tP05.0  = % 6.3fms\n"+
			"\tP25.0  = % 6.3fms\n"+
			"\tP50.0  = % 6.3fms\n"+
			"\tP75.0  = % 6.3fms\n"+
			"\tP95.0  = % 6.3fms\n"+
			"\tP99.5  = % 6.3fms\n"+
			"\tAVG    = % 6.3fms\n"+
			"\tDNSSEC = % 6t\n",
			res.TimeResults.P0Dot5.Seconds()*1000,
			res.TimeResults.P5.Seconds()*1000,
			res.TimeResults.P25.Seconds()*1000,
			res.TimeResults.P50.Seconds()*1000,
			res.TimeResults.P75.Seconds()*1000,
			res.TimeResults.P95.Seconds()*1000,
			res.TimeResults.P99Dot5.Seconds()*1000,
			res.TimeResults.Average.Seconds()*1000,
			res.DNSSECSupport,
		)
		fmt.Println("")
	}
}
