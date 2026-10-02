// securechan-curl — GET a path from a server over the PQ-hybrid compressed channel.
//
//	securechan-curl -addr host:8443 -pin <base64 server key> -path /api/v1/emily/collections/golden/items
//
// The bearer token (if any) is read from $TOKEN, never from argv (visible in `ps`).
package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"emily-agent/securechan"
)

func main() {
	addr := flag.String("addr", "", "host:port of a securechan listener")
	pin := flag.String("pin", "", "base64 pinned server public key")
	path := flag.String("path", "/", "request path")
	flag.Parse()
	key, err := base64.StdEncoding.DecodeString(*pin)
	if err != nil || *addr == "" {
		fmt.Fprintln(os.Stderr, "usage: securechan-curl -addr host:port -pin <base64> [-path /x]  (token via $TOKEN)")
		os.Exit(2)
	}
	hc := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, n, a string) (net.Conn, error) {
		return securechan.Dial(ctx, n, a, key)
	}}}
	req, _ := http.NewRequest("GET", "http://"+*addr+*path, nil)
	if t := os.Getenv("TOKEN"); t != "" {
		req.Header.Set("Authorization", "Bearer "+t)
	}
	resp, err := hc.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Fprintln(os.Stderr, resp.Status)
	io.Copy(os.Stdout, resp.Body)
}
