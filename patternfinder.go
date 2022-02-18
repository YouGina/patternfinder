package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type options struct {
	pattern     string
	depth       int
	concurrency int
	secure      bool
	debug       bool
}

func init() {
	flag.Usage = func() {
		h := []string{
			"",
			"Find patterns in http output based on regex string. Display occurences.",
			"",
			"Usage:",
			"    patternfinder [options] < urls.txt",
			"",
			"Options:",
			"    -p,\t\t--pattern <string> 	 Pattern to search for, default \"plugins/([[a-zA-Z0-9-_]+)/\"",
			"    -d,\t\t--depth <int>	  	 Depth to crawl, default 1",
			"    -c,\t\t--concurrency <int>	 Concurrency Level, default 2",
			"    -s,\t\t--secure  		 Enable TLS verification, default false",
			"    -dbg,\t--debug 		 Print all found patterns for debugging, default false",
			"",
		}
		fmt.Fprint(os.Stderr, strings.Join(h, "\n"))
	}
}

func main() {
	opts := options{}

	flag.StringVar(&opts.pattern, "pattern", "plugins/([[a-zA-Z0-9-_]+)/", "")
	flag.StringVar(&opts.pattern, "p", "plugins/([[a-zA-Z0-9-_]+)/", "")

	flag.IntVar(&opts.depth, "depth", 1, "")
	flag.IntVar(&opts.depth, "d", 1, "")

	flag.IntVar(&opts.concurrency, "concurrency", 2, "")
	flag.IntVar(&opts.concurrency, "c", 2, "")

	flag.BoolVar(&opts.secure, "secure", false, "")
	flag.BoolVar(&opts.secure, "s", false, "")

	flag.BoolVar(&opts.debug, "debug", false, "")
	flag.BoolVar(&opts.debug, "dbg", false, "")

	flag.Parse()

	if opts.pattern == "" {
		fmt.Fprintln(os.Stderr, "Pattern required. Hint: patternfinder --pattern \"plugins/([[a-zA-Z0-9-_]+)/\" < urls.txt")
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No urls detected. Hint: patternfinder < urls.txt")
		os.Exit(1)
	}

	cReg := regexp.MustCompile(opts.pattern)
	results := make(chan string, *&opts.concurrency)

	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			urlString := s.Text()

			u, err := url.Parse(urlString)
			if err != nil {
				log.Println("Error parsing URL:", err)
				return
			}
			hostname := u.Hostname()

			c := colly.NewCollector(
				colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0"),
				colly.AllowedDomains(hostname),
				colly.MaxDepth(*&opts.depth),
				colly.Async(true),
			)
			c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: opts.concurrency})
			c.WithTransport(&http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !*&opts.secure},
			})
			c.OnResponse(func(r *colly.Response) {

				body := string(r.Body)
				var sm sync.Map
				matches := cReg.FindAllStringSubmatch("`"+body+"`gm", -1)
				for _, v := range matches {
					_, present := sm.Load(v[1])
					if !present {
						sm.Store(v[1], true)
						if opts.debug {
							fmt.Fprintln(os.Stderr, v[1])
						}
					}
				}

				sm.Range(func(key, value interface{}) bool {
					results <- fmt.Sprint(key)
					return true
				})
			})

			c.Visit(urlString)
			c.Wait()
		}
		close(results)
	}()

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	patternResults := make(map[string]int)
	for res := range results {
		if _, ok := patternResults[res]; ok {
			patternResults[res] = patternResults[res] + 1
		} else {
			patternResults[res] = 1
		}
	}

	for key, value := range patternResults {
		fmt.Printf("%s\t%d\n", key, value)
	}
}
