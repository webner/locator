package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", printIP)

	fmt.Println("Hello World")

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "5000"
	}

	fmt.Printf("listening on %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func printIP(res http.ResponseWriter, req *http.Request) {

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	l := Location{ip, ""}

	// try proxy friendly headers
	for _, header := range []string{"X-Real-Ip", "X-Forwarded-For"} {
		realIP, ok := req.Header[header]
		if ok {
			fmt.Fprintln(res, realIP[0])
			return
		}
	}

	l.Office = getOffice(net.ParseIP(ip))

	b, err := json.Marshal(l)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fall back to the remote address
	res.Header().Set("Content-Type", "application/json")
	res.Write(b)
}

func getOffice(ip net.IP) string {

	mappings := []OfficeMapping{
		{"81.10.172.114/32", "Linz"},
		{"127.0.0.1/32", "localhost"},
	}

	for _, m := range mappings {
		_, addr, err := net.ParseCIDR(m.CIDR)
		if err == nil && addr.Contains(ip) {
			return m.Office
		}
	}

	return ""
}

type OfficeMapping struct {
	CIDR   string
	Office string
}

type Location struct {
	Ip     string
	Office string
}
