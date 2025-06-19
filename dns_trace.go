package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
)

func handleTrace(c *gin.Context) {
	_ = TraceFunc(c.Request.Context(), "handleTrace", func(ctx context.Context) error {
		domain := c.Query("domain")
		if domain == "" {
			c.String(http.StatusBadRequest, "Missing domain")
			return nil
		}

		if !strings.HasSuffix(domain, ".") {
			domain += "."
		}

		hops, err := traceDNS(ctx, domain)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return nil
		}

		c.String(http.StatusOK, strings.Join(hops, "\n"))
		return nil
	})
}

func traceDNS(ctx context.Context, domain string) ([]string, error) {
	var hops []string
	err := TraceFunc(ctx, "traceDNS", func(ctx context.Context) error {
		server := "198.41.0.4:53" // a.root-servers.net
		name := domain

		for {
			msg := new(dns.Msg)
			msg.SetQuestion(name, dns.TypeNS)
			c := new(dns.Client)
			resp, _, err := c.Exchange(msg, server)
			if err != nil {
				return fmt.Errorf("query failed at %s: %v", server, err)
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

			hops = append(hops, "")
			nsHost := nsRecords[0]

			msgA := new(dns.Msg)
			msgA.SetQuestion(nsHost, dns.TypeA)
			resolver := "8.8.8.8:53"
			respA, _, err := c.Exchange(msgA, resolver)
			if err != nil || len(respA.Answer) == 0 {
				break
			}

			found := false
			for _, a := range respA.Answer {
				if aRecord, ok := a.(*dns.A); ok {
					server = fmt.Sprintf("%s:53", aRecord.A.String())
					found = true
					break
				}
			}
			if !found {
				break
			}
		}

		hops = append(hops, "")
		// Final resolution
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

		return nil
	})
	return hops, err
}
