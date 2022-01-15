package udefiner

import (
	"fmt"
	"responder/go-app/internal/term"
	"strings"
)

type service struct {
	repo Repo
}

func (s service) GetUrbanDefinerDefinition(text string) (string, error) {
	t := term.Term(text)

	definitions, err := s.repo.GetUrbanDictionaryDefinition(t)

	if err != nil {
		return "", err
	}

	return formatDefinitions(t.Raw(), definitions), nil
}

func formatDefinitions(term string, uds []UrbanDefinition) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*%v*\n\n", term)

	for _, ud := range uds {
		fmt.Fprintf(&sb, "%v\n\n", ud)
	}

	return sb.String()
}

func newService() service {
	return service{newRepo()}
}
