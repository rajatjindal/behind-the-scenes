package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/variables"
	"github.com/rajatjindal/behind-the-scenes/api/pkg/slack"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
)

type Handler struct {
	slack *slack.Client

	triggerEmojiCode   string
	allowedChannelCode string
}

type Option func(h *Handler)

func WithTriggerEmoji(emoji string) Option {
	return func(h *Handler) {
		h.triggerEmojiCode = emoji
	}
}

func WithAllowedChannel(channel string) Option {
	return func(h *Handler) {
		h.allowedChannelCode = channel
	}
}

func NewHandler() (*Handler, error) {
	client := spinhttp.NewClient()

	sclient, err := slack.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to init slack client %w", err)
	}

	allowedChannel, err := variables.Get("allowed_channel")
	if err != nil {
		return nil, fmt.Errorf("allowed channel is required %w", err)
	}

	triggerEmoji, err := variables.Get("trigger_on_emoji_code")
	if err != nil {
		return nil, fmt.Errorf("trigger emoji is required %w", err)
	}

	h := &Handler{
		slack: sclient,
	}

	for _, option := range []Option{
		WithAllowedChannel(allowedChannel), WithTriggerEmoji(triggerEmoji),
	} {
		option(h)
	}

	if h.allowedChannelCode == "" {
		return nil, fmt.Errorf("allowed channel config is mandatory")
	}

	if h.triggerEmojiCode == "" {
		return nil, fmt.Errorf("trigger-emoji config is mandatory")
	}

	return h, nil
}

func (s *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	logrus.Info(string(dump))
	logrus.Info(r.Header)
	ctx := r.Context()

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	logrus.Infof("raw msg is %s", string(raw))
	err = s.slack.VerifySignature(r.Header, raw)
	if err != nil {
		logrus.Errorf("error when verifying signature %v", err)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	logrus.Info("signature verification successful")

	outerEvent, err := slack.ParseEvent(raw)
	if err != nil {
		logrus.Errorf("failed to parse slack event. err: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	logrus.Info("slack.ParseEvent successful")

	switch {
	case outerEvent.Type == slackevents.URLVerification:
		logrus.Info("handling url verification request")
		s.urlVerificationHandler(w, raw)
		return
	case outerEvent.Type == slackevents.CallbackEvent:
		err = s.handleCallbackEvent(ctx, outerEvent)
		if err != nil {
			logrus.Errorf("failed to handle event. err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	logrus.Infof("unknown event type %v", outerEvent.Type)
	fmt.Fprintln(w, "OK")
}

func (s *Handler) urlVerificationHandler(w http.ResponseWriter, raw []byte) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(raw), &r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text")
	w.Write([]byte(r.Challenge))
}
