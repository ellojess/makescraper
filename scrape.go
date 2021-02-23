package main

import (
	"fmt"
	"regexp"

	"encoding/csv"
	// "encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// AmazonResult struct that can be used for json
type AmazonResult struct {
	Content []string `json:"content"`
}

func readFile(file string) string {

	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(fileContents)
}

func writeFile(file string, data string) {
	bytesToWrite := []byte(data)
	err := ioutil.WriteFile(file, bytesToWrite, 0644)

	if err != nil {
		panic(err)
	}
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

func writeToJSON() {
	// fileName := "output.json"
	// file, err := os.Create(fileName)
	// if err != nil {
	// 	log.Fatalf("Could not create %s", fileName)
	// }

	// file, _ := json.MarshalIndent(data, "", " ")

	// _ = ioutil.WriteFile("output.json", file, 0644)

	// if err != nil {
	//   log.Fatalf("Could not create %s", fileName)
	// }
}

func writeToCSV() {
	fileName := "output.csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Could not create %s", fileName)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Product Name", "Stars", "Price"})

}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector; gives access to methods allowing
	// trigger callback functions when certain event happens
	c := colly.NewCollector(
		colly.Async(true),
	)

	// set random delays between every request
	c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second,
		Parallelism: 4, // max number of requests to be executed at a time
	})

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
	// c.Visit("https://www.amazon.com/s?k=alpaca+plush&ref=nb_sb_noss_1")
	// scrape data from 20 pages on Amazon result page
	for i := 1; i <= 20; i++ {
		fullURL := fmt.Sprintf("https://www.amazon.com/s?k=alpaca+plush&page=%d", i)
		c.Visit(fullURL)
	}
	c.Wait() // wait until all concurrent requests are done

	writeToCSV()
}
