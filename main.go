package main

import (
	"log"
	"net/http"
	"os"

	// Plug in CoreDNS
	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
)

func init() {
	dnsserver.Directives = append(dnsserver.Directives, "leak")
}

func main() {
	// Running without args just runs the DNS server
	if len(os.Args) == 1 {
		coremain.Run()
		return
	}

	switch os.Args[1] {
	case "serve":
		// Otherwise we run the HTTP server
		log.Fatal(http.ListenAndServe(":3333", newRouter()))
		return
	}

	log.Fatalf("unknown command %q", os.Args[1])
}
