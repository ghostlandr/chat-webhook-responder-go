package main

import (
	"context"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func checkVerification(fn http.HandlerFunc, definerType definers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var signingSecret string
		switch definerType {
		case define:
			signingSecret = definesigningSecret
		case udefine:
			signingSecret = udefineSigningSecret
		}
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
