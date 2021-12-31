package definer

import "responder/go-app/internal/term"

type service struct {
	repo repo
}

func (s service) GetDefinerDefinition(text string) (string, error) {
	t := term.Term(text)

	// Make request to dictionary site
	return s.repo.GetDictionaryDefinition(t)
}

func newService() service {
	return service{repo: newRepo()}
}
