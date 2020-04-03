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

type germanGrammarRule struct {
	Nominativ string
	Akkusativ string
	Dativ     string
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
						return identifyWord(noun.De, c.Args().Get(0))
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
						return identifyWord(noun.En, c.Args().Get(0))
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
	grammerRules := make(map[string]germanGrammarRule)
	grammerRules["Der"] = germanGrammarRule{"der / ein", "den / einen", "dem / einem"}
	grammerRules["Die"] = germanGrammarRule{"die / eine", "die / eine", "die / einer"}
	grammerRules["Das"] = germanGrammarRule{"das / ein", "das / ein", "dem / einem"}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"German", "English", "Nominativ", "Akusativ", "Dativ"})
	for _, n := range nouns {
		rule := grammerRules[n.Article]
		row := []string{n.Article + " " + n.De, n.En, rule.Nominativ, rule.Akkusativ, rule.Dativ}
		table.Append(row)
	}
	table.Render()
	return nil
}

func identifyWord(noun string, arg string) bool {
	return strings.ToLower(noun) == strings.ToLower(arg)
	//return strings.Contains(strings.ToLower(noun), strings.ToLower(arg))
}
