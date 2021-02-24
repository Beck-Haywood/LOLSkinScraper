package main
// Questions:
// Best way to define champion list
// how to loop through more then 1 scrape
// how to add to the file more then one struct of structs
import (
        "fmt"
        "strings"
        "os"
        "regexp"
        "encoding/json"
        "io/ioutil"
		"github.com/gocolly/colly"
)
type Skin struct {
        Name string `json: "name"`
        Cost string `json: "cost"`
		Date string `json: "date"`
    }

type SkinsStruct struct {
        Skins []Skin
    }

func main() {
		var names []string
		var costs []string
		var dates []string

		// champions := []string{"Aatrox", "Ahri", "Camille"}
		// championMap := make(map[string][]string)
		// make() creates the data in memory given the definition

        // Instantiate default collector
        c := colly.NewCollector()
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

        // Start scraping on https://leagueoflegends.fandom.com/wiki/Camille/LoL/Cosmetics
        c.Visit("https://leagueoflegends.fandom.com/wiki/Camille/LoL/Cosmetics")

        // fmt.Printf("%v", names)
        // fmt.Printf("%v", costs)
        skins := []Skin{}
        skinsStruct := SkinsStruct{skins}

        for i, name := range names {
                item1 := Skin{Name: name, Cost: costs[i], Date: dates[i]}
                skinsStruct.AddItem(item1)
                fmt.Println(name, costs[i], dates[i])
        }

        fmt.Printf("%+v\n", skinsStruct) // Print Struct with Variable Name
        b, err := json.Marshal(skinsStruct)
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Stdout.Write(b)
		_ = ioutil.WriteFile("output.json", b, 0644)
}

func (skins *SkinsStruct) AddItem(skin Skin) []Skin {
        skins.Skins = append(skins.Skins, skin)
        return skins.Skins
    }