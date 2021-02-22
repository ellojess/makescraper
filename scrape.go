package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
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

	extensions.RandomUserAgent(c) // have Colly generate new User Agent string before every request

	// On every a element which has href attribute call callback
	// parse HTML
	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		// loop through products in the search result list
		e.ForEach("div.a-section.a-spacing-medium", func(_ int, e *colly.HTMLElement) {
			// access wanted values with css selectors
			var productName, stars, price string

			productName = e.ChildText("span.a-size-medium.a-color-base.a-text-normal")
			if productName == "" {
				// If no name then return and go directly to the next element
				return
			}

			// [REVIEW] inconsistency found when getting stars and price
			stars = e.ChildText("span.a-icon-alt")
			formatStars(&stars)

			price = e.ChildText("span.a-price > span.a-offscreen")
			formatPrice(&price)

			fmt.Printf("Product Name: %s \nStars: %s \nPrice: %s \n", productName, stars, price)
		})
	})

	// Before making a request print "Visiting ..."
	// request Amazon's result page
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		fmt.Println("UserAgent", r.Headers.Get("User-Agent"))
	})

	// Start scraping on https://www.amazon.com
	c.Visit("https://www.amazon.com/s?k=alpaca+plush&ref=nb_sb_noss_1")
}
