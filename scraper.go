package main

import (
	"io/ioutil"
    "fmt"
	"github.com/gocolly/colly"
	"strings"
	"regexp"
	"encoding/json"
	// "os"
)

type Characters struct {
	CharacterMap map[string]Character `json: "Character"`
}
type Character struct {
	Name string  `json: "name"`
	Skins []Skin `json: "skins"`
}
type Skin struct {
	Name string `json: "name"`
	Cost string `json: "cost"`
	Date string `json: "date"`
}

func main() {
	charStruc := linkScrape()

	b, err := json.MarshalIndent(charStruc, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	// os.Stdout.Write(b)
	_ = ioutil.WriteFile("output.json", b, 0644)
}

// scrapes all of the links for each champion
func linkScrape() Characters{
	c := colly.NewCollector()
	var characters Characters
	characters.CharacterMap = make(map[string]Character) 

	var numOfLinks int
    linkSelector := ".label-only a"
	links := make(chan string, 180)
	result := make(chan Character, 1)

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
	go worker(links, result)
	go worker(links, result)
	go worker(links, result)
	go worker(links, result)

	close(links)

	for i := 0; i < numOfLinks; i++ {
		struc := <- result
		
		
		characters.CharacterMap[struc.Name] = struc
	}

	// Marshal map to json

	// Two choices 1. Add json to json then write to file
	// create struct and add maps together

	return characters
}

// worker gorountine function that scrapes each website concurrently
func worker(links chan string, result chan Character) {
	for link := range links {
		result <- champScrape(link)
	}
}

func champScrape(link string) Character {
	fmt.Println(link)
	var names []string
	var costs []string
	var dates []string
	var name string

	c := colly.NewCollector()
	c.OnHTML(".mw-redirect", func(e *colly.HTMLElement) {
		name = e.Text
	})
	c.OnHTML(".skin-icon+ div div:nth-child(1)", func(e *colly.HTMLElement) {
			name := strings.Replace(e.Text, " View in 3D", "", 1)
			names = append(names, name)
	})
	c.OnHTML(".skin-icon+ div div+ div", func(e *colly.HTMLElement) {
			m1 := regexp.MustCompile(`^[^/]+`)
			m2 := regexp.MustCompile(`[^/]*$`)

			res := m1.FindString(e.Text)
			cost := strings.TrimSpace(res)

			res2 := m2.FindString(e.Text)
			date := strings.TrimSpace(res2)

			costs = append(costs, cost)
			dates = append(dates, date)
	})
	

	c.OnError(func(_ *colly.Response, err error) {
			fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
			fmt.Println("Visited", r.Request.URL)
	})
	// c.OnScraped(func(r *colly.Response) {
	//         fmt.Println("Finished", r.Request.URL)
	// })
	c.Visit(link)
	// fmt.Println(name)
    // listOfChampSkins := ListofSkinsStruct{} // Init struct for each champ skinStruct

	skins := []Skin{} // Init list skins
	// skinsStruct := SkinsStruct{skins} // Init struct of skins


	for i, name := range names {
		item1 := Skin{Name: name, Cost: costs[i], Date: dates[i]} // Add all data using index
		// skinsStruct.AddItem(item1)
		skins = append(skins, item1)
		// fmt.Println(name, costs[i], dates[i])
	}
	// listOfChampSkins.ListofChampionSkins = append(listOfChampSkins.ListofChampionSkins, skinsStruct) // Append champion skins to listOfChampSkins struct
	var char Character
	char.Name = name 
	char.Skins = skins
	
	return char

}
