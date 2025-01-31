package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"
)

var (
	userAgent       = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36`
	site            = `https://book.codefine.site:6870/`
	cookie          = ""
	dir             = "./"
	timeout         = time.Duration(10) * time.Second
	username        = ""
	password        = ""
	verbose         = false
	startIndex      = 0
	version         = false
	continueOnStart = false
	retry           = 3
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
	flag.BoolVar(&version, "version", false, "show progream version")
	flag.BoolVar(&continueOnStart, "continue", true, "continue an incomplete download")

	flag.IntVar(&startIndex, "start-index", startIndex, "start book id")
	flag.IntVar(&retry, "retry", retry, "timeout retries")

	flag.Parse()

	if version {
		PrintVersion()
		os.Exit(0)
	}
}
func main() {
	tale, err := NewTableBook(site,
		WithRetry(retry),
		WithVerboseOption(verbose),
		WithUserCookieOption(cookie),
		WithUserAgentOption(userAgent),
		WithTimeOutOption(timeout),
		WithStartIndex(startIndex),
		WithLoginOption(username, password),
		WithContinue(continueOnStart),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d books retrieved on server %s", tale.Total, site)

	for {
		book, err := tale.Next()
		if err != nil {
			if errors.Is(err, NO_MORE_BOOK_ERROR) {
				log.Printf("%s [exit]", err.Error())
				if tale.exit != nil {
					tale.exit()
					os.Exit(0)
				}
			}
			log.Printf("%s [skiped]", err.Error())
			continue
		}

		if err = tale.Download(book, dir); err != nil {
			log.Printf("[%d/%d] downloading %s, %s [skiped]", book.Book.ID, tale.LastIndex(), book.Book.Title, err)
			continue
		}
		log.Printf("[%d/%d] downloading %s successed", book.Book.ID, tale.LastIndex(), book.String())
	}
}
