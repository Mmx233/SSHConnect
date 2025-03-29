package main

import (
	"flag"
	"fmt"
	"github.com/Mmx233/SSHConnect/tools"
	"log"
	"net/url"
	"os"
	"sync"

	"golang.org/x/net/proxy"
)

func main() {
	var socks5Addr = flag.String("S", "", "socks5 address")
	flag.Parse()

	socks5Url, err := url.Parse("//" + *socks5Addr)
	if err != nil {
		log.Fatalln("invalid socks5 address:", err)
	}

	if len(flag.Args()) < 2 {
		log.Fatalln("Usage: connect -S user@proxy:port <host> <port>")
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

	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%s", flag.Args()[0], flag.Args()[1]))
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
