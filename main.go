// A basic Go web server that traces DNS resolution steps using miekg/dns and Gin framework

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	httpsRouter := httpsRouter()
	httpRedirect := proxyRouter()

	// HTTP and HTTPS servers
	httpSrv := &http.Server{Addr: ":80", Handler: httpRedirect}
	httpsSrv := &http.Server{
		Addr:    ":443",
		Handler: httpsRouter,
	}
	pprofSrv := &http.Server{
		Addr:    ":6060",
		Handler: http.DefaultServeMux, // pprof handlers
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

	go func() {
		log.Println("📊 pprof listening on :6060")
		if err := pprofSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("pprof server failed: %v", err)
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
	if err := pprofSrv.Shutdown(ctx); err != nil {
		log.Fatalf("profiler Shutdown: %v", err)
	}

	fmt.Println("✅ Graceful shutdown complete")
}

func httpsRouter() *gin.Engine {
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

	// This middleware creates root spans for HTTP requests
	httpsRouter.Use(otelgin.Middleware("dns_say_what"))
	return httpsRouter
}

func proxyRouter() *gin.Engine {
	// tls/ssl redirect all traffic from :80 to :443
	httpRedirect := gin.Default()
	httpRedirect.Use(func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]
		target := "https://" + host + c.Request.RequestURI
		c.Redirect(http.StatusMovedPermanently, target)
	})
	return httpRedirect
}
