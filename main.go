// A basic Go web server that traces DNS resolution steps using miekg/dns and Gin framework

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
)

func main() {
	// Gin routers, serve two functionality:
	// 1. serve static react build
	// 2. accept /trace GET request
	httpsRouter := gin.Default()
	httpsRouter.Use(cors.Default())
	httpsRouter.Static("/static", "./client/build/static")
	httpsRouter.StaticFile("/favicon.ico", "./client/build/favicon.ico")
	httpsRouter.StaticFile("/manifest.json", "./client/build/manifest.json")
	httpsRouter.GET("/trace", handleTrace)
	httpsRouter.NoRoute(func(c *gin.Context) {
		c.File("./client/build/index.html")
	})

	// tls/ssl redirect all traffic from :80 to :443
	httpRedirect := gin.Default()
	httpRedirect.Use(func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]
		target := "https://" + host + c.Request.RequestURI
		c.Redirect(http.StatusMovedPermanently, target)
	})

	// HTTP and HTTPS servers
	httpSrv := &http.Server{Addr: ":80", Handler: httpRedirect}
	httpsSrv := &http.Server{
		Addr:    ":443",
		Handler: httpsRouter,
	}

	// Run both in goroutines
	go func() {
		fmt.Println("🚀 HTTPS server on https://localhost")
		if err := httpsSrv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	go func() {
		fmt.Println("🌐 HTTP redirect server on http://localhost")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n🛑 Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpsSrv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTPS Shutdown: %v", err)
	}
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP Shutdown: %v", err)
	}

	fmt.Println("✅ Graceful shutdown complete")
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

		hops = append(hops, "")
		// Resolve one of the next NS records to IP
		nsHost := nsRecords[0]
		msgA := new(dns.Msg)
		msgA.SetQuestion(nsHost, dns.TypeA)
		resolver := "8.8.8.8:53" // instead of using 'server'
		respA, _, err := c.Exchange(msgA, resolver)
		if err != nil || len(respA.Answer) == 0 {
			break
		}

		for _, a := range respA.Answer {
			if aRecord, ok := a.(*dns.A); ok {
				server = fmt.Sprintf("%s:53", aRecord.A.String())
				break
			}
		}
	}

	hops = append(hops, "")
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
