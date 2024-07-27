package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

type server struct {
	address string
	query   string
	record  dns.Type
}

var servers = []server{
	{address: "resolver1.opendns.com:53", query: "myip.opendns.com", record: dns.Type(dns.TypeA)},
	{address: "resolver2.opendns.com:53", query: "myip.opendns.com", record: dns.Type(dns.TypeA)},
	{address: "resolver3.opendns.com:53", query: "myip.opendns.com", record: dns.Type(dns.TypeA)},
	{address: "resolver4.opendns.com:53", query: "myip.opendns.com", record: dns.Type(dns.TypeA)},
	{address: "ns1.google.com:53", query: "o-o.myaddr.l.google.com", record: dns.Type(dns.TypeTXT)},
	{address: "ns2.google.com:53", query: "o-o.myaddr.l.google.com", record: dns.Type(dns.TypeTXT)},
	{address: "ns3.google.com:53", query: "o-o.myaddr.l.google.com", record: dns.Type(dns.TypeTXT)},
	{address: "ns4.google.com:53", query: "o-o.myaddr.l.google.com", record: dns.Type(dns.TypeTXT)},
}

func resolve(s server) string {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(s.query), uint16(s.record))

	c := new(dns.Client)
	c.Timeout = time.Second * 3
	resp, _, err := c.Exchange(m, s.address)
	if err != nil {
		fmt.Printf("Failed to query DNS server: %v\n", err)
		return ""
	}

	if len(resp.Answer) == 0 {
		return ""
	}

	for _, answer := range resp.Answer {
		switch uint16(s.record) {
		case uint16(dns.TypeTXT):
			if txtRecord, ok := answer.(*dns.TXT); ok {
				fmt.Println(s.address, txtRecord.Txt[0])
				return txtRecord.Txt[0]
			}
		case uint16(dns.TypeA):
			if aRecord, ok := answer.(*dns.A); ok {
				fmt.Println(s.address, aRecord.A.String())
				return aRecord.A.String()
			}
		}
	}

	return ""
}

func shuffle(s []server) []server {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	serversShuffled := shuffle(servers)

	for _, s := range serversShuffled {
		myIp := resolve(s)
		if myIp != "" {
			fmt.Fprintln(w, myIp)
			return
		}
	}

	fmt.Fprintln(w, "-1")
}

func main() {
	http.HandleFunc("/", rootHandler)
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
