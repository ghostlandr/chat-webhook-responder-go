package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/logging"
	"github.com/anaskhan96/soup"
	"github.com/slack-go/slack"

	"responder/go-app/internal/definer"
)

var (
	slackTokens = []string{
		"70XqnEL12zOlA08Fo0lraciE",
		"ayYWtEzhfqh5GcXdEqrD3H3h",
		"7b5WbqiybRqPDRTm2e9GvTUL",
		"x56o3ZQzti2l7YEb7ntRu4gE",
	}
	definitionLimit = 5
)

func isAuthorizedToken(token string) bool {
	for _, tok := range slackTokens {
		if tok == token {
			return true
		}
	}
	return false
}

type Term string

func (t Term) Raw() string {
	return strings.Split(string(t), ": ")[1]
}

func (t Term) String() string {
	term := t.Raw()
	st := strings.Split(strings.ToLower(term), " ")
	return strings.Join(st, "-")
}

type slackResponse struct {
	Text string `json:"text"`
}

func udefineHandler(w http.ResponseWriter, r *http.Request) {
	// http://www.urbandictionary.com/define.php?term=
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests supported", 405)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("something is bad with your data: %s", err), 400)
		return
	}

	token := r.FormValue("token")
	text := r.FormValue("text")

	if !isAuthorizedToken(token) {
		http.Error(w, "not authorized", 403)
		return
	}
	t := Term(text)
	log.Printf("Going to search for %s\n", t)
	dictURL := fmt.Sprintf("https://www.urbandictionary.com/define.php?term=%s", t)
	log.Printf("Going to search for %s\n", dictURL)
	resp, err := soup.Get(dictURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s\n", err), 400)
		return
	}
	log.Printf("%v", resp)
	doc := soup.HTMLParse(resp)
	log.Printf("%v", doc)
	main := doc.Find("div", "class", "row")
	if main.Error != nil {
		renderResponse(w, fmt.Sprintf("No results for %s", t.Raw()))
		return
	}

	content := main.Find("div", "id", "content")
	if content.Error != nil {
		renderResponse(w, fmt.Sprintf("No results for %s", t.Raw()))
		return
	}

	definitions := content.FindAll("div", "class", "def-panel")

	for _, defPanel := range definitions {
		ud := NewUrbanDefinition(defPanel)
		fmt.Printf("%v", ud)
	}
}

func hackIt() {
	t := Term("define: butt")

	dictURL := fmt.Sprintf("https://www.urbandictionary.com/define.php?term=%s", t)

	resp, err := soup.Get(dictURL)
	if err != nil {
		return
	}

	doc := soup.HTMLParse(resp)

	// main := doc.Find("div", "class", "row")
	// if main.Error != nil {
	// 	log.Printf("No results for %s 1", t.Raw())
	// 	return
	// }

	content := doc.Find("div", "id", "content")
	if content.Error != nil {
		log.Printf("No results for %s 2", t.Raw())
		return
	}

	definitions := content.FindAll("div", "class", "def-panel")

	for _, defPanel := range definitions[:1] {
		_ = NewUrbanDefinition(defPanel)
	}
}

type UrbanDefinition struct {
	Definition string
	Example    string
	Upvotes    uint32
	Downvotes  uint32
}

// New expects HTML that includes the .def-panel class and its contents
func NewUrbanDefinition(defPanel soup.Root) UrbanDefinition {
	if defPanel.Error != nil {
		return UrbanDefinition{}
	}

	definition := defPanel.Find("div", "class", "meaning").FullText()

	example := defPanel.Find("div", "class", "example").FullText()

	thumbsArea := defPanel.Find("div", "class", "thumbs")

	log.Printf("%#v", thumbsArea.Find("a", "class", "up").Find("span", "class", "count").Pointer)
	log.Printf("%v", thumbsArea.Find("a", "class", "up").FullText())

	upvotes := thumbsArea.Find("a", "class", "up").Find("span", "class", "count").Text()
	downvotes := thumbsArea.Find("a", "class", "down").Find("span", "class", "count").Text()
	log.Printf("Up: %v, down: %v", upvotes, downvotes)

	return UrbanDefinition{
		Definition: definition, Example: example,
	}
}

func defineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests supported", 405)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("something is bad with your data: %s", err), 400)
		return
	}

	token := r.FormValue("token")
	text := r.FormValue("text")

	if !isAuthorizedToken(token) {
		http.Error(w, "not authorized", 403)
		return
	}
	t := Term(text)

	// Make request to dictionary site
	log.Printf("Going to search for %s\n", t)
	dictURL := fmt.Sprintf("https://en.oxforddictionaries.com/definition/%s", t)
	resp, err := soup.Get(dictURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s\n", err), 400)
		return
	}
	doc := soup.HTMLParse(resp)
	main := doc.Find("section", "class", "gramb")
	if main.Error != nil {
		renderResponse(w, fmt.Sprintf("No results for %s", t.Raw()))
		return
	}
	defs := main.FindAll("span", "class", "ind")
	grammar := main.Find("span", "class", "pos")

	var o string
	for i, d := range defs {
		if i < definitionLimit {
			o += fmt.Sprintf("%d. %s\n", i+1, d.Text())
		}
	}
	if len(o) > 0 {
		o = fmt.Sprintf("Definitions for *%s* - _%s_\n\n", t.Raw(), grammar.Text()) + o
		o += fmt.Sprintf("\n_Brought to you by <%s|English Oxford Dictionaries>_", dictURL)
	}

	renderResponse(w, o)
}

func renderResponse(w http.ResponseWriter, o string) {
	err := json.NewEncoder(w).Encode(slackResponse{Text: o})
	if err != nil {
		http.Error(w, fmt.Sprintf("not sure what happened: %s", err), 400)
	}
}

func main() {
	ctx := context.Background()
	projectID := "chat-webhook-responder"
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	logger := client.Logger("something").StandardLogger(logging.Info)

	mux := http.NewServeMux()
	mux.HandleFunc("/define", checkVerification(definer.New(logger).ServeDefinerRequest))
	// http.HandleFunc("/define/slack", defineHandler)
	// http.HandleFunc("/udefine", udefineHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Printf("Error? %s", http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func checkVerification(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			log.Printf("verification failed, but letting things through for now in testing phase, %v", err)
			// w.WriteHeader(http.StatusUnauthorized)
			// return
		}

		// If all is well, pass the request along
		fn(w, r)
	}
}

func addContext(ctx context.Context, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(ctx)
		fn(w, req)
	}
}
