package tokens

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"responder/go-app/internal/config"
	"responder/go-app/internal/logs"
	"strings"
)

type Definers int64

const (
	Undefined Definers = iota
	Define
	Udefine
)

var slackAuthURL string = "https://slack.com/oauth/v2/authorize?scope=%s&client_id=%s"

type oauthHandler struct {
	logger logs.Logger
}

type OauthHandler interface {
	ServeDefinerOauthAuthorizeRequest(w http.ResponseWriter, r *http.Request)
	ServeUrbanDefinerOauthAuthorizeRequest(w http.ResponseWriter, r *http.Request)
}

func NewOauthHandler(l logs.Logger) OauthHandler {
	return oauthHandler{l}
}

func (o oauthHandler) ServeDefinerOauthAuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	scopes := []string{"commands"}
	clientID := config.DefinerClientID
	clientSecret := config.DefinerClientSecret

	o.handleActualRequest(strings.Join(scopes, ","), clientID, clientSecret, "Definin'", w, r)
}

func (o oauthHandler) ServeUrbanDefinerOauthAuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	scopes := []string{"commands"}
	clientID := config.UdefinerClientID
	clientSecret := config.UdefinerClientSecret

	o.handleActualRequest(strings.Join(scopes, ","), clientID, clientSecret, "Udefinin'", w, r)
}

type OauthV2AccessResponse struct {
	OK          bool   `json:"ok"`
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

func (o oauthHandler) handleActualRequest(scopes, clientID, clientSecret, app string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		o.logger.Printf("error parsing form, somehow: %s", err)
		return
	}

	code, gotCode := r.Form["code"]
	o.logger.Printf("got a code: %v", code)

	if !gotCode {
		redirectURL := fmt.Sprintf(slackAuthURL, scopes, clientID)
		o.logger.Printf("redirecting to %s", redirectURL)

		http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
		return
	}

	o.logger.Printf("have a code so making a call to exchange tokens")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code[0])

	o.logger.Printf("requesting access token from Slack")

	resp, err := http.Post("https://slack.com/api/oauth.v2.access", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		o.logger.Printf("error getting access token: %s", err)
		http.Error(w, fmt.Sprintf("error getting access token: %s", err), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	var oauthResponse OauthV2AccessResponse
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		o.logger.Printf("error reading response from Slack: %s", err)
		return
	}

	if err = json.Unmarshal(bodyBytes, &oauthResponse); err != nil {
		o.logger.Printf("error unmarshalling json: %s", err)
		return
	}

	if !oauthResponse.OK {
		o.logger.Printf("error connecting new workspace: %s", oauthResponse.Error)
		http.Error(w, fmt.Sprintf("Something went wrong getting you connected, try again later? %s", oauthResponse.Error), http.StatusBadRequest)
		return
	}

	o.logger.Printf("connected new workspace successfully")
	fmt.Fprintf(w, "Succesfully connected! Now start %s!", app)
}
