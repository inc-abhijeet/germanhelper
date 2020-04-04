package germanhelper

import (
	"fmt"
	"strings"

	"github.com/cameronnorman/arcee"
	"github.com/gocolly/colly"
)

func SearchGermanWord(word string) (result []GermanNoun) {
	findFunc := func(noun GermanNoun) bool {
		return wordMatcher(noun.De, word)
	}
	result = searchCache(findFunc, word)

	if len(result) < 1 {
		result, _ = searchWebsiteDeEn(word)
	}

	return result
}

func SearchEnglishWord(word string) (result []GermanNoun) {
	findFunc := func(noun GermanNoun) bool {
		return wordMatcher(noun.En, word)
	}
	result = searchCache(findFunc, word)

	if len(result) < 1 {
		result, _ = searchWebsiteEnDe(word)
	}

	return result
}

func searchCache(findFunc func(noun GermanNoun) bool, word string) (result []GermanNoun) {
	germanNouns := fetchCached()
	arcee.Select(findFunc, germanNouns, &result)

	return result
}

func wordMatcher(noun string, arg string) bool {
	return strings.ToLower(noun) == strings.ToLower(arg)
	//return strings.Contains(strings.ToLower(noun), strings.ToLower(arg))
}

func searchWebsiteDeEn(arg string) (nouns []GermanNoun, err error) {
	var articleIdentifier = "html body.src-de.target-en div#main-grid section#content div#inner-content article.lemma-group.src-de.target-en"
	var wordIdentifier = "html body.src-de.target-en div#main-grid section#content div#inner-content article.lemma-group.src-de.target-en div.tab div div.tab-content div.tab-inner-content div.summary div.search-term h5"
	var genderIdentifier = "html body.src-de.target-en div#main-grid section#content div#inner-content article.lemma-group.src-de.target-en div.tab div div.tab-content div.tab-inner-content div.summary div.search-term span.dict-additions span.abbr abbr"
	var translationIdentifier = "html body.src-de.target-en div#main-grid section#content div#inner-content article.lemma-group.src-de.target-en div.tab div div.tab-content div.tab-inner-content div.summary div.summary-inner a.btn.blue.round span.btn-inner"
	genderMapping := make(map[string]string)
	genderMapping["m"] = "Der"
	genderMapping["f"] = "Die"
	genderMapping["n"] = "Das"
	genderMapping["m pl"] = "Die (Plural)"
	genderMapping["f pl"] = "Die (Plural)"
	genderMapping["n pl"] = "Die (Plural)"

	c := colly.NewCollector()
	c.OnHTML(articleIdentifier, func(e *colly.HTMLElement) {
		var siteGender = e.ChildText(genderIdentifier)
		fmt.Println(siteGender)
		var siteWord = e.ChildText(wordIdentifier)
		var siteTrans = e.ChildText(translationIdentifier)

		if siteTrans != "" {
			nouns = append(nouns, GermanNoun{siteTrans, siteWord, genderMapping[siteGender], ""})
		}
	})

	c.Visit("https://en.langenscheidt.com/german-english/" + arg)

	return nouns, err
}

func searchWebsiteEnDe(arg string) (nouns []GermanNoun, err error) {
	var articleIdentifier = "html body.src-en.target-de div#main-grid section#content div#inner-content article.lemma-group.src-en.target-de"
	var wordIdentifier = "html body.src-en.target-de div#main-grid section#content div#inner-content article.lemma-group.src-en.target-de div.tab div label.lemma-pos-title.flag.en h2"

	genderMapping := make(map[string]string)
	genderMapping["m"] = "Der"
	genderMapping["f"] = "Die"
	genderMapping["n"] = "Das"
	genderMapping["m pl"] = "Die (Plural)"

	c := colly.NewCollector()
	c.OnHTML(articleIdentifier, func(e *colly.HTMLElement) {
		var siteWord = e.ChildText(wordIdentifier)
		if strings.Contains(siteWord, "noun") == true {
			var translationsIdentifier = "html body.src-en.target-de div#main-grid section#content div#inner-content article.lemma-group.src-en.target-de div.tab div div.tab-content div.tab-inner-content div.senses div.sense-item"
			var genderIdentifier = "html body.src-en.target-de div#main-grid section#content div#inner-content article.lemma-group.src-en.target-de div.tab div div.tab-content div.tab-inner-content div.senses div.sense-item ul.lemma-entry-group li.lemma-entry.translation div.col1 div.trans-line div.inter span.lemma-pieces span.trans span.pos span.abbr abbr"
			e.ForEach(translationsIdentifier, func(_ int, el *colly.HTMLElement) {
				var translationIdentifier = "html body.src-en.target-de div#main-grid section#content div#inner-content article.lemma-group.src-en.target-de div.tab div div.tab-content div.tab-inner-content div.senses div.sense-item ul.lemma-entry-group li.lemma-entry.translation div.col1 div.trans-line div.inter span.lemma-pieces span.trans span.trans a"
				var siteTrans = el.ChildText(translationIdentifier)
				var siteGender = el.ChildText(genderIdentifier)
				if len(nouns) < 1 {
					nouns = append(nouns, GermanNoun{siteWord, siteTrans, genderMapping[siteGender], ""})
				}
			})
		}
	})

	c.Visit("https://en.langenscheidt.com/english-german/" + arg)

	return nouns, err
}
