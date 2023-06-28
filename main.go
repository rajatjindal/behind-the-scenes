package main

import (
	"io"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/pets-of-fermyon/pkg/slack"
	"github.com/slack-go/slack/slackevents"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		parentEvent, err := slack.ParseEvent(raw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if parentEvent.Type == slackevents.URLVerification {
			slack.URLVerificationHandler(w, raw)
			return
		}

	})
}

func main() {}
