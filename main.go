// A basic Go web server that traces DNS resolution steps using miekg/dns and Gin framework
// Run: go run main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
)

func main() {
	r := gin.Default()

	r.GET("/", serveForm)
	r.GET("/trace", handleTrace)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

func serveForm(c *gin.Context) {
	html := `<html>
		<head><title>DNS Trace</title></head>
		<body>
		<h1>DNS Trace Tool</h1>
		<form action="/trace">
		  Domain: <input name="domain" type="text" />
		  <input type="submit" value="Trace DNS" />
		</form>
		</body>
		</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func handleTrace(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.String(http.StatusBadRequest, "Missing domain")
		return
	}

	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	hops, err := traceDNS(domain)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, strings.Join(hops, "\n"))
}

func traceDNS(domain string) ([]string, error) {
	hops := []string{}
	server := "198.41.0.4:53" // a.root-servers.net
	name := domain

	for {
		msg := new(dns.Msg)
		msg.SetQuestion(name, dns.TypeNS)
		c := new(dns.Client)
		resp, _, err := c.Exchange(msg, server)
		if err != nil {
			return hops, fmt.Errorf("query failed at %s: %v", server, err)
		}

		line := fmt.Sprintf("Queried %s → NS: ", server)
		nsRecords := []string{}
		for _, ans := range resp.Ns {
			if ns, ok := ans.(*dns.NS); ok {
				nsRecords = append(nsRecords, ns.Ns)
			}
		}
		line += strings.Join(nsRecords, ", ")
		hops = append(hops, line)

		if len(nsRecords) == 0 {
			break
		}

		// Resolve one of the next NS records to IP
		nsHost := nsRecords[0]
		msgA := new(dns.Msg)
		msgA.SetQuestion(nsHost, dns.TypeA)
		respA, _, err := c.Exchange(msgA, server)
		if err != nil || len(respA.Answer) == 0 {
			break // can't resolve next hop
		}

		for _, a := range respA.Answer {
			if aRecord, ok := a.(*dns.A); ok {
				server = fmt.Sprintf("%s:53", aRecord.A.String())
				break
			}
		}
	}

	// Final A/AAAA resolution
	msg := new(dns.Msg)
	msg.SetQuestion(domain, dns.TypeA)
	resp, _, err := new(dns.Client).Exchange(msg, server)
	if err == nil {
		for _, a := range resp.Answer {
			if aRec, ok := a.(*dns.A); ok {
				hops = append(hops, fmt.Sprintf("Final A record from %s: %s", server, aRec.A.String()))
			}
		}
	}

	return hops, nil
}
