package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RenderResponse takes in a ResponseWriter and a SlackResponse and writes the SlackResponse
// to the 'Writer.
func RenderResponse(w http.ResponseWriter, o SlackResponse) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(o)
	if err != nil {
		http.Error(w, fmt.Sprintf("not sure what happened: %s", err), 400)
	}
}

func RenderStringInChannel(w http.ResponseWriter, s string) {
	sr := New(s).InChannel()
	RenderResponse(w, sr)
}

// slackResponse{Text: o, ResponseType: "in_channel"}

type SlackResponse interface {
	InChannel() SlackResponse
}

func New(s string) SlackResponse {
	return &slackResponse{
		ResponseType: "ephemeral",
		Blocks: []block{
			{
				Type: "section",
				Text: textBlock{
					Type: "mrkdwn",
					Text: s,
				},
			},
		},
	}
}

func (s *slackResponse) InChannel() SlackResponse {
	s.ResponseType = "in_channel"
	return s
}

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
	BlockID string    `json:"block_id,omitempty"`
}

type textBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type slackResponse struct {
	Blocks       []block `json:"blocks,omitempty"`
	Text         string  `json:"text,omitempty"`
	ResponseType string  `json:"response_type,omitempty"`
}
