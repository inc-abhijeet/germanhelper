package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cameronnorman/arcee"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

type germanNoun struct {
	En      string `json:"en"`
	De      string `json:"de"`
	Article string `json:"article"`
	Plural  string `json:"plural"`
}

func main() {
	germanNouns := []germanNoun{}
	data, err := ioutil.ReadFile("words.json")
	if err != nil {
		log.Fatal("Unable to load German Nouns")
	}

	err = json.Unmarshal(data, &germanNouns)
	if err != nil {
		log.Fatal("Unable to parse German Nouns file into struct")
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "german noun",
				Aliases: []string{"de"},
				Usage:   "de Mug",
				Action: func(c *cli.Context) error {
					result := []germanNoun{}
					findFunc := func(noun germanNoun) bool {
						return strings.Contains(strings.ToLower(noun.De), strings.ToLower(c.Args().Get(0)))
					}
					arcee.Select(findFunc, germanNouns, &result)
					if len(result) > 0 {
						err = printResult(result)
						if err != nil {
							log.Fatal("Unable to print result")
						}
					}
					fmt.Println("Noun not found. If valid please inform the developer to add it to the list")

					return nil
				},
			},
			{
				Name:    "english noun",
				Aliases: []string{"en"},
				Usage:   "en Wort",
				Action: func(c *cli.Context) error {
					result := []germanNoun{}
					findFunc := func(noun germanNoun) bool {
						return strings.Contains(strings.ToLower(noun.En), strings.ToLower(c.Args().Get(0)))
					}
					arcee.Select(findFunc, germanNouns, &result)
					if len(result) > 0 {
						err = printResult(result)
						if err != nil {
							log.Fatal("Unable to print result")
						}
					}
					fmt.Println("Noun not found. If valid please inform the developer to add it to the list")

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func printResult(nouns []germanNoun) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Article", "German", "English"})
	for _, n := range nouns {
		row := []string{n.Article, n.Article + " " + n.De, n.En}
		table.Append(row)
	}
	table.Render()
	return nil
}
