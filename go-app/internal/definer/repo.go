package definer

import (
	"fmt"
	"log"
	"responder/go-app/internal/term"

	"github.com/anaskhan96/soup"
)

var definitionLimit = 5

type repo struct{}

func (r repo) GetDictionaryDefinition(term term.Term) (string, error) {
	log.Printf("Going to search for %s\n", term)
	dictURL := fmt.Sprintf("https://en.oxforddictionaries.com/definition/%s", term)
	resp, err := soup.Get(dictURL)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	doc := soup.HTMLParse(resp)
	main := doc.Find("section", "class", "gramb")
	if main.Error != nil {
		return "", fmt.Errorf("no results for %s", term.Raw())
	}
	defs := main.FindAll("span", "class", "ind")
	grammar := main.Find("span", "class", "pos")

	var o string
	for i, d := range defs {
		if i < definitionLimit {
			o += fmt.Sprintf("%d. %s\n", i+1, d.Text())
		}
	}
	if len(o) > 0 {
		o = fmt.Sprintf("Definitions for *%s* - _%s_\n\n", term.Raw(), grammar.Text()) + o
		o += fmt.Sprintf("\n_Brought to you by <%s|English Oxford Dictionaries>_", dictURL)
	}
	return o, nil
}

func newRepo() repo {
	return repo{}
}
