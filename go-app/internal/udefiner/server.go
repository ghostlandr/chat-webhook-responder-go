package udefiner

import (
	"fmt"
	"net/http"
	"responder/go-app/internal/logs"
	"responder/go-app/internal/response"
	"strings"

	"github.com/slack-go/slack"
)

type server struct {
	service Service
	logger  logs.Logger
}

func New(l logs.Logger) UrbanDefiner {
	return server{newService(), l}
}

func (s server) ServeUrbanDefinerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests supported", http.StatusMethodNotAllowed)
		return
	}

	command, err := slack.SlashCommandParse(r)
	if err != nil {
		s.logger.Printf("slack command parse failed: %v", err)
		http.Error(w, fmt.Sprintf("something is bad with your data: %s", err), http.StatusBadRequest)
		return
	}

	text := command.Text

	s.logger.Printf("Going to search for %v with token %v", text)

	result, err := s.service.GetUrbanDefinerDefinition(text)
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		s.logger.Printf(errStr)
		if strings.HasPrefix(errStr, "error") {
			http.Error(w, errStr, 400)
		} else if strings.HasPrefix(errStr, "no results") {
			response.RenderStringPrivately(w, fmt.Sprintf("No results found for %v", text))
		} else {
			response.RenderStringPrivately(w, fmt.Sprintf("Something went wrong trying to search for %v", text))
		}
		return
	}

	response.RenderStringInChannel(w, result)
}
