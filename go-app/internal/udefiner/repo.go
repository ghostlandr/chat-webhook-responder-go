package udefiner

import (
	"errors"
	"fmt"
	"responder/go-app/internal/term"

	"github.com/anaskhan96/soup"
)

type repo struct{}

func (r repo) GetUrbanDictionaryDefinition(t term.Term) ([]UrbanDefinition, error) {
	dictURL := fmt.Sprintf("https://www.urbandictionary.com/define.php?term=%s", t)

	resp, err := soup.Get(dictURL)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	doc := soup.HTMLParse(resp)

	uds, err := getDefinitionsFromHTML(doc)
	if err != nil {
		return nil, err
	}

	return uds, err
}

func getDefinitionsFromHTML(doc soup.Root) ([]UrbanDefinition, error) {
	definitions := doc.FindAll("div", "class", "definition")

	if len(definitions) == 0 {
		return nil, errors.New("no results found")
	}

	defIDs := make([]string, 0, len(definitions))
	defs := make([]UrbanDefinition, 0, len(definitions))
	for _, defPanel := range definitions[:5] {
		ud := NewUrbanDefinition(defPanel)

		defs = append(defs, ud)
		defIDs = append(defIDs, ud.DefID)
	}

	thumbsData, err := getVotesForDefinitionIDs(defIDs)
	if err != nil {
		return nil, err
	}

	// Double for-loop is scary, but the max amount of definitions we can have is 5.
	// If it was ever higher than 5 we could reconsider the data structures here.
	// For now this reads well and will run at most 25 times total - easy.
	for _, data := range thumbsData {
		for idx, def := range defs {
			if def.DefID == data.DefID() {
				def.SetVotes(data.Upvotes(), data.Downvotes())
				defs[idx] = def
				break
			}
		}
	}

	return defs, nil
}

func newRepo() repo {
	return repo{}
}
