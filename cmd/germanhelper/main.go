package main

import (
	"fmt"
	"os"

	"germanhelper"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type germanGrammarRule struct {
	Nominativ string
	Akkusativ string
	Dativ     string
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "german noun",
				Aliases: []string{"de"},
				Usage:   "de Wort",
				Action: func(c *cli.Context) error {
					result := germanhelper.SearchGermanWord(c.Args().Get(0))
					if len(result) < 1 {
						fmt.Println("Noun not found. If valid please inform the developer to add it to the list")
						return nil
					}

					printResult(result)
					return nil
				},
			},
			{
				Name:    "english noun",
				Aliases: []string{"en"},
				Usage:   "en roof",
				Action: func(c *cli.Context) error {
					result := germanhelper.SearchEnglishWord(c.Args().Get(0))
					if len(result) < 1 {
						fmt.Println("Noun not found. If valid please inform the developer to add it to the list")
						return nil
					}

					printResult(result)
					return nil
				},
			},
			{
				Name:    "pronomen",
				Aliases: []string{"pro"},
				Usage:   "pro",
				Action: func(c *cli.Context) error {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Nominativ", "Akkusativ", "Dativ", "Reflexiv", "Possessiv / Adjecktiv"})
					table.SetColumnColor(
						tablewriter.Colors{tablewriter.FgHiWhiteColor},
						tablewriter.Colors{tablewriter.FgHiWhiteColor},
						tablewriter.Colors{tablewriter.FgHiWhiteColor},
						tablewriter.Colors{tablewriter.FgHiWhiteColor},
						tablewriter.Colors{tablewriter.FgHiWhiteColor},
					)
					row := []string{"ich", "mich", "mir", "mich", "mein - "}
					table.Append(row)
					row = []string{"du", "dich", "dir", "dich", "dein - "}
					table.Append(row)
					row = []string{"er", "ihn", "ihm", "sich", "sein - "}
					table.Append(row)
					row = []string{"sie", "sie", "ihr", "sich", "ihr - "}
					table.Append(row)
					row = []string{"es", "es", "ihm", "sich", "sein - "}
					table.Append(row)
					row = []string{"wir", "uns", "uns", "uns", "unser - "}
					table.Append(row)
					row = []string{"ihr", "euch", "euch", "euch", "euer - "}
					table.Append(row)
					row = []string{"sie", "sie", "ihnen", "sich", "ihr - "}
					table.Append(row)
					row = []string{"Sie", "Sie", "ihnen", "Sich", "Ihr - "}
					table.Append(row)

					table.Render()
					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}

func printResult(nouns []germanhelper.GermanNoun) error {
	grammerRules := make(map[string]germanGrammarRule)
	grammerRules["Der"] = germanGrammarRule{"der / ein", "den / einen", "dem / einem"}
	grammerRules["Die"] = germanGrammarRule{"die / eine", "die / eine", "die / einer"}
	grammerRules["Das"] = germanGrammarRule{"das / ein", "das / ein", "dem / einem"}
	grammerRules["Die (Plural)"] = germanGrammarRule{"die", "die", "den"}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"German", "English", "Nominativ", "Akusativ", "Dativ"})
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)
	for _, n := range nouns {
		rule := grammerRules[n.Article]
		row := []string{n.Article + " " + n.De, n.En, rule.Nominativ, rule.Akkusativ, rule.Dativ}
		table.Append(row)
	}
	table.Render()
	return nil
}
