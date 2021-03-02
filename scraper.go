package main

import (
	"io/ioutil"
    "fmt"
	"github.com/gocolly/colly"
	"strings"
)

// scrapes all of the links for each champion
func linkScrape() {
	c := colly.NewCollector()

	var numOfLinks int
    linkSelector := ".label-only a"
	links := make(chan string, 180)
	result := make(chan string, 1)


	c.OnHTML(linkSelector, func(e *colly.HTMLElement) {
        // fmt.Println(e.Text)
		link := e.Attr("href")
        links <- ("https://leagueoflegends.fandom.com" + link + "/Cosmetics")
		numOfLinks ++
	})

	c.Visit("https://leagueoflegends.fandom.com/wiki/List_of_champions")

	go worker(links, result)
	go worker(links, result)
	go worker(links, result)
	go worker(links, result)
	go worker(links, result)

	close(links)

	for i := 0; i < numOfLinks; i++ {
		<- result
	}

	return
}

// writes a .json file that holds the json of the articles
func writeJSONFile(json []byte) {
	err := ioutil.WriteFile("output.json", json, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	linkScrape()
	// Test ChampScrape Works
	// champScrape("https://leagueoflegends.fandom.com/wiki/Azir/LoL/Cosmetics")
}

// worker gorountine function that scrapes each website concurrently
func worker(links chan string, result chan string) {
	for link := range links {
		result <- champScrape(link)
	}
}

func champScrape(link string) string{
	var name string
	fmt.Println(link)
	c := colly.NewCollector()
	c.OnHTML(".skin-icon+ div div:nth-child(1)", func(e *colly.HTMLElement) {
			name = strings.Replace(e.Text, " View in 3D", "", 1)
	})
	c.Visit(link)
	// fmt.Println(name)
	return name

}