package webhook

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rajatjindal/pets-of-fermyon/pkg/bluesky"
	"github.com/rajatjindal/pets-of-fermyon/pkg/slack"
	"github.com/slack-go/slack/slackevents"
)

type Handler struct {
	slack   *slack.Client
	bluesky *bluesky.BlueSky
}

func NewHandler(slack *slack.Client, bluesky *bluesky.BlueSky) *Handler {
	// api := slack.New("xoxb-4430321508054-5503179271856-u2b6oMyPEhqnZunK0PBy1e32", slacksdk.OptionHTTPClient(httpclient))

	return &Handler{
		slack:   slack,
		bluesky: bluesky,
	}
}

func (s *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outerEvent, err := slack.ParseEvent(raw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if outerEvent.Type == slackevents.URLVerification {
		slack.URLVerificationHandler(w, raw)
		return
	}

	if outerEvent.Type != slackevents.CallbackEvent {
		fmt.Fprintln(w, "OK")
		return
	}

	err = s.handleCallbackEvent(ctx, outerEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "OK")
}
