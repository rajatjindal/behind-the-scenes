package main

import (
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/pets-of-fermyon/pkg/bluesky"
	"github.com/rajatjindal/pets-of-fermyon/pkg/creds/kvcreds"
	"github.com/rajatjindal/pets-of-fermyon/pkg/slack"
	"github.com/rajatjindal/pets-of-fermyon/pkg/webhook"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		credsProvider := kvcreds.Provider()
		client := spinhttp.NewClient()

		bsky, err := bluesky.NewClient(client, credsProvider)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		sclient, err := slack.NewClient(client, credsProvider)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		handler := webhook.NewHandler(sclient, bsky)
		handler.Handle(w, r)
	})
}

func main() {}
