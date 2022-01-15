package udefiner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type GetVotesResponse struct {
	Thumbs []Thumbs `json:"thumbs"`
}

type Thumbs struct {
	DefIDInt     int    `json:"defid"`
	UpvotesInt   int    `json:"up"`
	DownvotesInt int    `json:"down"`
	Current      string `json:"current"`
}

func (t Thumbs) Upvotes() string {
	return strconv.Itoa(t.UpvotesInt)
}

func (t Thumbs) Downvotes() string {
	return strconv.Itoa(t.DownvotesInt)
}

func (t Thumbs) DefID() string {
	return strconv.Itoa(t.DefIDInt)
}

var getVotesURL string = "https://api.urbandictionary.com/v0/uncacheable?ids=%s"

func getVotesForDefinitionIDs(defIDs []string) ([]Thumbs, error) {
	requestURL := fmt.Sprintf(getVotesURL, strings.Join(defIDs, ","))
	resp, err := http.Get(requestURL)

	if err != nil {
		// TODO: handle
		return nil, err
	}
	defer resp.Body.Close()

	var voteResponse GetVotesResponse
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bodyBytes, &voteResponse); err != nil {
		return nil, err
	}

	return voteResponse.Thumbs, nil
}
