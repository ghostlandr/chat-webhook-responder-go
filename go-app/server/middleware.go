package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"responder/go-app/internal/config"
)

type definers int64

const (
	undefined definers = iota
	define
	udefine
)

func checkVerification(fn http.HandlerFunc, definerType definers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var signingSecret string
		switch definerType {
		case define:
			signingSecret = config.DefineSigningSecret
		case udefine:
			signingSecret = config.UdefineSigningSecret
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Reset the body so it can be used in future handlers
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		verifier, err := NewSecretsVerifier(r.Header, signingSecret, string(body[:]))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			// log.Printf("verification failed, but letting things through for now in testing phase, %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
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
