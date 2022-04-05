package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"golang.org/x/time/rate"
)

var (
	userAgent  = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36`
	site       = `https://book.codefine.site:6870/`
	cookie     = ""
	dir        = "./"
	timeout    = time.Duration(10) * time.Second
	username   = ""
	password   = ""
	concurrent = 1
	verbose    = false
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flag.StringVar(&cookie, "cookie", cookie, "http cookie")
	flag.StringVar(&username, "username", username, "username")
	flag.StringVar(&password, "password", password, "password")
	flag.StringVar(&site, "site", site, "tabebook web site")
	flag.StringVar(&dir, "dir", dir, "data dir")
	flag.StringVar(&userAgent, "user-agent", userAgent, "http userAgent")
	flag.DurationVar(&timeout, "timeout", timeout, "http timeout")
	flag.BoolVar(&verbose, "verbose", false, "show debug log")
	flag.IntVar(&concurrent, "c", concurrent, "maximum number of concurrent download tasks allowed per second")

	flag.Parse()
}
func main() {
	tale, err := NewTableBook(site,
		WithVerboseOption(verbose),
		WithUserCookieOption(cookie),
		WithUserAgentOption(userAgent),
		WithTimeOutOption(timeout),
		WithLoginOption(username, password),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d books retrieved on server %s", tale.ServerInfo.Sys.Books, site)
	l := rate.NewLimiter(rate.Limit(concurrent), concurrent)

	for {
		// 限制速度
		l.Wait(context.Background())
		book, err := tale.Next()
		if err != nil {
			log.Printf("%s [skiped]", err.Error())
			if errors.Is(err, NO_MORE_BOOK_ERROR) {
				os.Exit(0)
			}
			continue
		}

		go func() {
			if err = tale.Download(book, dir); err != nil {
				log.Printf("downloading %s , found %s [skiped]", book.Book.Title, err)
				return
			}
			log.Printf("downloading %s successed", book.String())
		}()
	}
}
