package germanhelper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
)

type GermanNoun struct {
	En      string `json:"en"`
	De      string `json:"de"`
	Article string `json:"article"`
	Plural  string `json:"plural"`
}

type GermanVerb struct {
	En string `json:"en"`
	De string `json:"de"`
}

func init() {
	fetchCached("nouns")
	fetchCached("verbs")
}

func fetchCached(wordType string) []GermanNoun {
	germanNouns := []GermanNoun{}
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(usr.HomeDir + "/config/germanhelper/" + wordType + ".json"); err == nil {
		// path/to/whatever exists
		data, err := ioutil.ReadFile("words.json")

		if err != nil {
			log.Fatal("Unable to load German Nouns")
		}

		err = json.Unmarshal(data, &germanNouns)
		if err != nil {
			log.Fatal("Unable to parse German Nouns file into struct")
		}

	} else if os.IsNotExist(err) {
		client := http.Client{}
		res, err := client.Get("https://raw.githubusercontent.com/cameronnorman/germanhelper/master/" + wordType + ".json")
		if err != nil {
			log.Fatal("Unable to fetch required data")
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Unable to parse response body")
		}
		err = saveOffline(wordType+".json", body)
		if err != nil {
			log.Fatal("Unable to save data")
		}
		err = json.Unmarshal(body, &germanNouns)
		if err != nil {
			log.Fatal("Unable to parse German Nouns file into struct")
		}
	}

	return germanNouns
}

func saveOffline(filename string, data []byte) error {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(usr.HomeDir+"/config/germanhelper/"+filename, data, 0644)
	return err
}

//func FetchVerbs(url string) {
//germanVerbs := []GermanVerb{}
//data, err := ioutil.ReadFile("verbs.json")

//if err != nil {
//log.Fatal("Unable to load German Verbs")
//}

//err = json.Unmarshal(data, &germanVerbs)
//if err != nil {
//log.Fatal("Unable to parse German Verbs file into struct")
//}

//c := colly.NewCollector()
//c.OnHTML("html body main div.wrap div.inner div.lang.lang-parse div.table1-box", func(e *colly.HTMLElement) {
//e.ForEach("table.table1 tr", func(_ int, el *colly.HTMLElement) {
//en := el.ChildText("td:first-child")
//de := el.ChildText("td:nth-child(2)")
//germanVerbs = append(germanVerbs, GermanVerb{en, de})
//})
//})

//c.Visit(url)

//germanVerbsJson, _ := json.Marshal(germanVerbs)
//err = ioutil.WriteFile("verbs.json", germanVerbsJson, 0644)
//return
//}
