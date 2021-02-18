package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"io/ioutil"
)


// AmazonResult struct that can be used for json
type AmazonResult struct {
	Content []string `json:"content"`
}

func formatPrice(price *string) {
	r := regexp.MustCompile(`\$(\d+(\.\d+)?).*$`)

	newPrices := r.FindStringSubmatch(*price)

	if len(newPrices) > 1 {
		*price = newPrices[1]
	} else {
		*price = "Unknown"
	}

}

func formatStars(stars *string) {
	if len(*stars) >= 3 {
		*stars = (*stars)[0:3]
	} else {
		*stars = "Unknown"
	}
}

// Writes data onto file
func writeFile(name string, data string) {
	bytesToWrite := []byte(data)
	err := ioutil.WriteFile(name, bytesToWrite, 0644)
	if err != nil {
		panic(err)
	}
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector; gives access to methods allowing 
	// trigger callback functions when certain event happens
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	// 
	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {
                link := e.Attr("href")

				// Print link
                fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://www.amazon.com/
	c.Visit("https://www.amazon.com/s?k=alpaca+plush&ref=nb_sb_noss_1")
}
