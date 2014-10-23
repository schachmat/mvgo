package main

import (
	"fmt"
	"strings"
//	"io/ioutil"
	"log"
	"os"
	"net/http"
	"net/url"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/charset"
)

const uri = "http://www.mvg-live.de/ims/dfiStaticAuswahl.svc?"
var params = []string{}
var gotstation = false

func parse_arg(t string) {
	if t[0] == '-' {
		for i := 1; i < len(t); i++ {
			if t[i] == 'b' {
				params = append(params, "bus=checked")
			} else if t[i] == 's' {
				params = append(params, "sbahn=checked")
			} else if t[i] == 't' {
				params = append(params, "tram=checked")
			} else if t[i] == 'u' {
				params = append(params, "ubahn=checked")
			}
		}
	} else if !gotstation {
		params = append(params, "haltestelle=" + url.QueryEscape(t))
		gotstation = true
	}
}

func parse_departure_line(n *html.Node) {
	var line, station, departure string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode || c.Data != "td" {
			continue
		}
		for _, a := range c.Attr {
			if a.Key != "class" {
				continue
			}
			if a.Val == "lineColumn" {
				line = c.FirstChild.Data
				break
			}
			if a.Val == "stationColumn" {
				station = strings.TrimSpace(c.FirstChild.Data)
				break
			}
			if a.Val == "inMinColumn" {
				departure = c.FirstChild.Data
				break
			}
		}
	}
	fmt.Println(line, station, departure)
}

func parse_departures(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		for _, a := range n.Attr {
			if a.Key == "class" && (a.Val == "rowOdd" || a.Val == "rowEven") {
				parse_departure_line(n)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse_departures(c)
	}
}

func parse_stations(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && strings.HasPrefix(a.Val, "/ims/dfiStatic") {
				fmt.Println(a.Val, n.FirstChild.Data)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse_stations(c)
	}
}

func find_table(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "departureTable departureView" {
				parse_departures(n)
				return n
			} else if a.Key == "class" && a.Val == "departureTable header" {
				parse_stations(n)
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		t := find_table(c)
		if t != nil {
			return t
		}
	}
	return nil
}

func main() {
	for i := 1; i < len(os.Args); i++ {
		parse_arg(os.Args[i])
	}

	res, err := http.Get(uri + strings.Join(params, "&"))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	r, err := charset.NewReader(res.Body, res.Header["Content-Type"][0])
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	find_table(doc)
}
