package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/labstack/echo/v4"
)

type categories struct {
	// Title is the title of the category
	// Link is the link of the category
	// Products is an array of all products in the category
	Title    string    `json:"categoryTitle"`
	Link     string    `json:"categoryLink"`
	// Products []product `json:"categoryProducts"` 
}

type product struct {
	// Name is the name of the product
	// Image is the image link associated with the product
	// Price is the price of the product
	// Link is the link to the product
	Name  string `json:"productName"`
	Image string `json:"productImage"`
	Price string `json:"productPrice"`
	Link string `json:"productLink"`
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector and echo object
	e := echo.New()
	c := colly.NewCollector(colly.Async(true))

	c.Limit(&colly.LimitRule{
		DomainGlob: "costco.com/*",
		RandomDelay: 2 * time.Second,
		Parallelism: 2,
	})

	// Get the Category
	categorySelector := "#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div"
	var categoriesList []categories

	c.OnHTML(categorySelector, func(e *colly.HTMLElement) {
		e.ForEach("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div", func(num int, h *colly.HTMLElement) {
			// var p []product
			fmt.Println("This iteration:", num)
			currentCount := strconv.Itoa(num + 1)

			categoryName := string("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div:nth-child(" + currentCount + ") > a")
			categoryLink := string("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div:nth-child(" + currentCount + ") > a")
			// fmt.Println(categoryName, "\n", categoryLink)

			cName := e.ChildText(categoryName)
			cLink := e.ChildAttr(categoryLink, "href")
			fmt.Println(cName,"\n",cLink)

			d := categories{Title: cName, Link: cLink}
			categoriesList = append(categoriesList, d)
			// fmt.Println(categoriesList)
		})
	})

	// c.OnHTML("" , func(b *colly.HTMLElement) {
	// 	b.ForEach("#search-results > ctl:cache > div.product-list.grid", func(_ int, g *colly.HTMLElement) {
	// 		productName := g.ChildText("#search-results > > div.product-list.grid > div > div > div.thumbnail > div.caption.link-behavior > div.caption > p.description > a")
	// 		productPrice := g.ChildText("#search-results > > div.product-list.grid > div > div > div.thumbnail > div.caption.link-behavior > div.caption > div > div")
	// 		productImage := g.ChildText("#search-results > > div.product-list.grid > div > div > div.thumbnail > div.product-img-holder.link-behavior > div > img")
	// 		productLink := g.ChildAttr("#search-results > > div.product-list.grid > div > div > div.thumbnail > div.caption.link-behavior > div.caption > p.description > a", "#search-results > > div.product-list.grid > div > div > div.thumbnail > div.caption.link-behavior > div.caption > p.description")
	// 		fmt.Println("Second HTML")

	// 		Adding individual products to the product list
	// 		pl := product{Name: productName, Price: productPrice, Image: productImage, Link: productLink}
	// 		p = append(p, pl)
	// 	})	
	// })
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping Costco
	// c.Visit("https://www.costco.com/all-costco-grocery.html")
	c.Visit("https://www.bjs.com/content?template=B&espot_main=EverydayEssentials&source=megamenu")


	fmt.Println(categoriesList)

	// Serve to echo
	e.GET("/scrape", func(f echo.Context) error {
		return f.JSON(http.StatusOK, categoriesList)
	})

	// Handle errors
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// After data is scraped, marshall to JSON
	DataJSONarr, err := json.MarshalIndent(categoriesList, "", "	")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("output.json", DataJSONarr, 0644)
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8000"))

	c.Wait()
}
