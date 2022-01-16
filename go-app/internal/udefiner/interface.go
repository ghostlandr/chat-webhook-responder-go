package udefiner

import (
	"net/http"
	"responder/go-app/internal/term"
)

type Repo interface {
	GetUrbanDictionaryDefinition(t term.Term) ([]UrbanDefinition, error)
}

type Service interface {
	GetUrbanDefinerDefinition(text string) (string, error)
}

type UrbanDefiner interface {
	ServeUrbanDefinerRequest(w http.ResponseWriter, r *http.Request)
}

type Logger interface {
	Printf(format string, v ...interface{})
}
