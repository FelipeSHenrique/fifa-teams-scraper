package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type Team struct {
	Team   string
	League string
	Ata    string
	Mei    string
	Def    string
	Ger    string
}

func writeCSV(filename string, data [][]string) error {
	csvFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	for _, v := range data {
		if err := csvWriter.Write(v); err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}
	}
	csvWriter.Flush()
	return nil
}

func main() {
	var teams []Team

	pageToScrape := "https://www.fifaindex.com/teams/?page=1"
	nextPage := ""

	i := 1
	limit := 23

	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.OnHTML("ul.pagination li.ml-auto a.btn", func(e *colly.HTMLElement) {
		attr := e.Attr("href")
		nextPage = fmt.Sprintf("https://www.fifaindex.com%s", attr)
	})

	c.OnHTML("table.table-teams tbody tr", func(e *colly.HTMLElement) {
		team := Team{}

		teamName := e.ChildText("td a.link-team")

		if teamName != "" {
			team.Team = teamName
			team.League = e.ChildText("td a.link-league")
			team.Ata = e.ChildText("td[data-title=ATT] span")
			team.Mei = e.ChildText("td[data-title=MID] span")
			team.Def = e.ChildText("td[data-title=DEF] span")
			team.Ger = e.ChildText("td[data-title=OVR] span")
			teams = append(teams, team)
		}
	})

	c.OnScraped(func(response *colly.Response) {
		if i < limit {
			i++
			c.Visit(nextPage)
		}
	})

	err := c.Visit(pageToScrape)
	if err != nil {
		log.Fatalln("Failed to visit site", err)
	}

	var res [][]string

	res = append(res, []string{
		"rank",
		"team",
		"league",
		"ata",
		"mei",
		"def",
		"ger",
	})

	for i, t := range teams {
		res = append(res, []string{
			fmt.Sprint(i + 1),
			t.Team,
			t.League,
			t.Ata,
			t.Mei,
			t.Def,
			t.Ger,
		})
	}

	err = writeCSV("rank-teams.csv", res)
	if err != nil {
		log.Fatalln("Failed to create CSV file", err)
	}

}
