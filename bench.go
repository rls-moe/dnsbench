package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type benchResult struct {
	DNSSECSupport bool
	TimeResults   TimeResults
}

type TimeResults struct {
	P0Dot5  time.Duration
	P5      time.Duration
	P25     time.Duration
	P50     time.Duration
	P75     time.Duration
	P95     time.Duration
	P99Dot5 time.Duration
	Average time.Duration
}

type ProgressCallback func(i, n uint16)

func bench(dnsServer, target string, measurements uint16, cb ProgressCallback) (*benchResult, error) {
	if !strings.Contains(dnsServer, ":") {
		dnsServer += ":53"
	}
	c := new(dns.Client)
	c.SingleInflight = true

	result := &benchResult{
		DNSSECSupport: false,
	}

	fmt.Println("Checking DNSSEC...")
	// verify DNSSEC
	m := new(dns.Msg)
	m.SetEdns0(4096, true) // Set DNSSEC OK
	m.SetQuestion("www.dnssec-failed.org.", dns.TypeA)
	r, _, err := c.Exchange(m, dnsServer)
	if err != nil {
		return nil, errors.Wrap(err, "Could not check DNSSEC")
	}
	if r.Rcode == dns.RcodeServerFailure {
		result.DNSSECSupport = true
	}

	// execute measurements
	var ttls = make([]time.Duration, 0)
	for i := measurements; i > 0; i-- {
		cb(measurements-i, measurements)
		time.Sleep(time.Duration(*sleepTimeout) * time.Millisecond)
		q := new(dns.Msg)
		if *antiCache {
			q.SetQuestion(RandStringRunes(4)+"."+target+".", dns.TypeA)
		} else {
			q.SetQuestion(target+".", dns.TypeA)
		}
		r, ttl, err := c.Exchange(q, dnsServer)
		if err != nil {
			return nil, errors.Wrap(err, "Bench Question failed")
		}
		if len(r.Answer) != 1 {
			return nil, errors.New("DNS has no response Answers")
		}
		ttls = append(ttls, ttl)
	}

	sort.Slice(ttls, func(i, j int) bool {
		// sort worst first
		return ttls[i].Nanoseconds() > ttls[j].Nanoseconds()
	})

	var avg int64
	for k := range ttls {
		avg += ttls[k].Nanoseconds()
	}
	avg /= int64(len(ttls))

	result.TimeResults.Average = time.Duration(avg) * time.Nanosecond

	result.TimeResults.P0Dot5 = Px(ttls, 0.005)
	result.TimeResults.P5 = Px(ttls, 0.05)
	result.TimeResults.P25 = Px(ttls, 0.25)
	result.TimeResults.P50 = Px(ttls, 0.50)
	result.TimeResults.P75 = Px(ttls, 0.75)
	result.TimeResults.P95 = Px(ttls, 0.95)
	result.TimeResults.P99Dot5 = Px(ttls, 0.995)

	return result, nil
}

func Px(sl []time.Duration, p float64) time.Duration {
	var res int64
	numR := int64(
		math.Min(
			math.Ceil(p*float64(len(sl))),
			float64(len(sl)),
		),
	)
	for i := int64(0); i < numR && i < int64(len(sl)); i++ {
		res += sl[i].Nanoseconds()
	}
	res /= numR
	return time.Duration(res) * time.Nanosecond
}
