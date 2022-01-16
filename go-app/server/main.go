package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"responder/go-app/internal/definer"
	"responder/go-app/internal/udefiner"

	"cloud.google.com/go/logging"
)

type definers int64

const (
	undefined definers = iota
	define
	udefine
)

func main() {
	ctx := context.Background()
	projectID := "chat-webhook-responder"
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	logger := client.Logger("chat-webhook-responder-go").StandardLogger(logging.Info)

	mux := http.NewServeMux()
	mux.HandleFunc("/define", checkVerification(definer.New(logger).ServeDefinerRequest, define))
	mux.HandleFunc("/udefine", checkVerification(udefiner.New(logger).ServeUrbanDefinerRequest, udefine))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Printf("Error? %s", http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func checkVerification(fn http.HandlerFunc, definerType definers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var signingSecret string
		// switch definerType {
		// case define:
		// 	signingSecret = definesigningSecret
		// case udefine:
		// 	signingSecret = udefineSigningSecret
		// }
		// verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		// if err = verifier.Ensure(); err != nil {
		// 	log.Printf("verification failed, but letting things through for now in testing phase, %v", err)
		// 	// w.WriteHeader(http.StatusUnauthorized)
		// 	// return
		// }

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
