package webhook

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/fermyon/spin/sdk/go/config"
	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/behind-the-scenes/pkg/creds/kvcreds"
	"github.com/rajatjindal/behind-the-scenes/pkg/logrus"
	"github.com/rajatjindal/behind-the-scenes/pkg/slack"
	"github.com/rajatjindal/behind-the-scenes/pkg/socialmedia"
	"github.com/slack-go/slack/slackevents"
)

type Handler struct {
	slack   *slack.Client
	socials []socialmedia.Provider

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
	credsProvider := kvcreds.Provider()
	client := spinhttp.NewClient()

	sclient, err := slack.NewClient(client, credsProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to init slack client %w", err)
	}

	allowedChannel, err := config.Get("allowed_channel")
	if err != nil {
		return nil, fmt.Errorf("allowed channel is required %w", err)
	}

	triggerEmoji, err := config.Get("trigger_on_emoji_code")
	if err != nil {
		return nil, fmt.Errorf("trigger emoji is required %w", err)
	}

	h := &Handler{
		slack:   sclient,
		socials: []socialmedia.Provider{},
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

	if outerEvent.Type == slackevents.URLVerification {
		logrus.Info("handling url verification request")
		slack.URLVerificationHandler(w, raw)
		return
	}
	logrus.Info("not a url verification request")

	if outerEvent.Type != slackevents.CallbackEvent {
		logrus.Info("not a callback event")
		fmt.Fprintln(w, "OK")
		return
	}
	logrus.Info("it is a callback request")

	err = s.handleCallbackEvent(ctx, outerEvent)
	if err != nil {
		logrus.Errorf("failed to handle event. err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info("request handled successfully")

	fmt.Fprintln(w, "OK")
}
