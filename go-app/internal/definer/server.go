package definer

import (
	"fmt"
	"net/http"
	"responder/go-app/internal/logs"
	"responder/go-app/internal/response"
	"strings"

	"github.com/slack-go/slack"
)

type DefinerServer interface {
	ServeDefinerRequest(w http.ResponseWriter, r *http.Request)
}

type server struct {
	logger  logs.Logger
	service service
}

func New(l logs.Logger) DefinerServer {
	return server{logger: l, service: newService()}
}

func (s server) ServeDefinerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests supported", http.StatusMethodNotAllowed)
		return
	}

	// err := r.ParseForm()
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("something is bad with your data: %s", err), http.StatusBadRequest)
		return
	}

	s.logger.Printf("receiving parameters: %v, %v", command.Token, command.Text)

	o, err := s.service.GetDefinerDefinition(command.Text)
	if err != nil {
		// log
		errStr := fmt.Sprintf("%v", err)
		s.logger.Printf("error handling define request: %s", errStr)
		if strings.HasPrefix(errStr, "error") {
			http.Error(w, errStr, 400)
		} else if strings.HasPrefix(errStr, "no results") {
			response.RenderStringPrivately(w, errStr)
		}
		return
	}

	response.RenderStringInChannel(w, o)
}
