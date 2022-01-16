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
	mux.HandleFunc("/define", addContext(ctx, checkVerification(definer.New(logger).ServeDefinerRequest, define)))
	mux.HandleFunc("/udefine", addContext(ctx, checkVerification(udefiner.New(logger).ServeUrbanDefinerRequest, udefine)))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Printf("Error? %s", http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
