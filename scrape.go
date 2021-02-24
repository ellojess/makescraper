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
	// ProductName string `json:"productName"`
	// Stars string `json:"stars"`
	// Price string `json:"price"`

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

	// hold data thats been scraped
	// var dataSlice []string

	// On every a element which has href attribute call callback
	// parse HTML

	// whole page
	// #search > div.s-desktop-width-max.s-desktop-content.sg-row > div.sg-col-16-of-20.sg-col.sg-col-8-of-12.sg-col-12-of-16 > div > span:nth-child(4) > div.s-main-slot.s-result-list.s-search-results.sg-row >

	// one item
	// #search > div.s-desktop-width-max.s-desktop-content.sg-row > div.sg-col-16-of-20.sg-col.sg-col-8-of-12.sg-col-12-of-16 > div > span:nth-child(4) > div.s-main-slot.s-result-list.s-search-results.sg-row > div:nth-child(13) > div > span > div > div
	// #search > div.s-desktop-width-max.s-desktop-content.sg-row > div.sg-col-16-of-20.sg-col.sg-col-8-of-12.sg-col-12-of-16 > div > span:nth-child(4) > div.s-main-slot.s-result-list.s-search-results.sg-row > div:nth-child(13) > div > span > div > div

	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {
		// c.OnHTML("div.s-result-item.s-search-results.sg-col-inner", func(e *colly.HTMLElement) {

		// loop through products in the search result list
		e.ForEach("div.a-section.a-spacing-medium", func(_ int, e *colly.HTMLElement) {
			// e.ForEach("div.s-result-item.s-search-results.sg-col-inner", func(_ int, e *colly.HTMLElement) {
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

			// print data being scraped
			fmt.Printf("Product Name: %s \nStars: %s \nPrice: %s \n", productName, stars, price)

			// dataSlice = append(dataSlice, e.Text)
			// data := AmazonResult{Content: dataSlice}
			// fmt.Println(data)

			// scrapedJSON, _ := json.MarshalIndent(data, "", "    ")
			// fmt.Println(string(scrapedJSON))

			// // write to JSON file
			// _ = ioutil.WriteFile("output.json", scrapedJSON, 0644)

			// // to append to a file
			// // create the file if it doesn't exists with O_CREATE, Set the file up for read write,
			// // add the append flag and set the permission
			// f, err := os.OpenFile("/var/log/debug-web.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// // write to file, f.Write()
			// f.Write(scrapedJSON)

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
	// writeToJSON()
}
