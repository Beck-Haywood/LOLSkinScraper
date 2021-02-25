package main
// Questions:
// Best way to define champion list
// how to loop through more then 1 scrape
// how to add to the file more then one struct of structs
import (
        "fmt"
        // "strings"
        // "os"
        // "regexp"
        // "encoding/json"
        // "io/ioutil"
        "github.com/gocolly/colly"
        "github.com/gocolly/colly/queue"

)

func main() {
		// Instantiate default collector
        c := colly.NewCollector()
        // create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
        c.OnHTML(".label-only a", func(e *colly.HTMLElement) {
			link := e.Attr("href")
                        q.AddURL("https://leagueoflegends.fandom.com" + link + "/Cosmetics")
                        // fmt.Printf("https://leagueoflegends.fandom.com" + link + "/Cosmetics")
			e.Request.Visit(link)
			scrapeSite(link)
        })
        c.OnRequest(func(r *colly.Request) {
                fmt.Println("Visiting", r.URL)
        })

        c.OnError(func(_ *colly.Response, err error) {
                fmt.Println("Something went wrong:", err)
        })

        c.OnResponse(func(r *colly.Response) {
                fmt.Println("Visited", r.Request.URL)
        })

        c.OnScraped(func(r *colly.Response) {
                fmt.Println("Finished", r.Request.URL)
        })

        c.Visit("https://leagueoflegends.fandom.com/wiki/List_of_champions")
        // fmt.Printf(q)
        q.Run(c)
}

func scrapeSite(link string) {
	fmt.Printf("Test1")
}