package main

import (
	"flag"
	"fmt"
	"github.com/Mmx233/SSHConnect/tools"
	"log"
	"net"
	"net/url"
	"os"
	"sync"

	"golang.org/x/net/proxy"
)

func main() {
	var socks5Addr = flag.String("S", "", "socks5 address")
	var resolve = flag.Bool("R", false, "resolve host to ip")
	flag.Parse()
	if len(flag.Args()) < 2 || *socks5Addr == "" {
		log.Fatalln("Usage: connect -S user@proxy:port <host> <port>")
	}

	host := flag.Args()[0]
	port := flag.Args()[1]

	if *resolve {
		result, err := net.LookupHost(host)
		if err != nil {
			log.Fatalln("resolve host failed:", err)
		} else if len(result) == 0 {
			log.Fatalln("resolve host failed: no address found")
		}
		host = result[0]
	}

	socks5Url, err := url.Parse("//" + *socks5Addr)
	if err != nil {
		log.Fatalln("invalid socks5 address:", err)
	}

	var auth *proxy.Auth
	if socks5Url.User != nil {
		auth = &proxy.Auth{User: socks5Url.User.Username()}
		pass, ok := socks5Url.User.Password()
		if !ok {
			auth.Password = pass
		}
	}

	dialer, err := proxy.SOCKS5("tcp", socks5Url.Host, auth, proxy.Direct)
	if err != nil {
		log.Fatalln("dial socks5 failed:", err)
	}

	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalln("dial failed:", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go tools.IOCopy(conn, os.Stdin, wg)
	go tools.IOCopy(os.Stdout, conn, wg)

	wg.Wait()
	_ = conn.Close()
}
