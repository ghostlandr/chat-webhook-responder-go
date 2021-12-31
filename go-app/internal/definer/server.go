package definer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"responder/go-app/tokens"
	"strings"

	"github.com/slack-go/slack"
)

type DefinerServer interface {
	ServeDefinerRequest(w http.ResponseWriter, r *http.Request)
}

type server struct {
	logger  Logger
	service service
}

type Logger interface {
	Printf(format string, v ...interface{})
}

func New(l Logger) DefinerServer {
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

	if !tokens.IsAuthorizedToken(command.Token) {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	o, err := s.service.GetDefinerDefinition(command.Text)
	if err != nil {
		// log
		errStr := fmt.Sprintf("%v", err)
		if strings.HasPrefix(errStr, "error") {
			http.Error(w, errStr, 400)
		} else if strings.HasPrefix(errStr, "no results") {
			renderResponse(w, errStr)
		}
		return
	}

	// fmt.Fprintf(w, o)
	renderResponse(w, o)
}

func renderResponse(w http.ResponseWriter, o string) {
	err := json.NewEncoder(w).Encode(inChannelMarkdown(o))
	if err != nil {
		http.Error(w, fmt.Sprintf("not sure what happened: %s", err), 400)
	}
}

// slackResponse{Text: o, ResponseType: "in_channel"}

func inChannelMarkdown(o string) interface{} {
	return slackResponse{
		ResponseType: "in_channel",
		Blocks: []block{
			{
				Type: "section",
				Text: textBlock{
					Type: "mrkdwn",
					Text: o,
				},
				BlockID: "definition",
			},
		},
	}
}

type block struct {
	Type    string    `json:"type"`
	Text    textBlock `json:"text"`
	BlockID string    `json:"block_id"`
}

type textBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type slackResponse struct {
	Blocks       []block `json:"blocks,omitempty"`
	Text         string  `json:"text,omitempty"`
	ResponseType string  `json:"response_type"`
}
