package slack

import (
	"encoding/json"
	"net/http"

	"github.com/fermyon/spin-go-sdk/variables"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Client struct {
	signingSecret string
	*slack.Client
}

func NewClient(httpclient *http.Client) (*Client, error) {
	token, err := variables.Get("slack_token")
	if err != nil {
		return nil, err
	}

	signingSecret, err := variables.Get("slack_signing_secret")
	if err != nil {
		return nil, err
	}

	client := slack.New(token, slack.OptionHTTPClient(httpclient))
	return &Client{
		Client:        client,
		signingSecret: signingSecret,
	}, nil
}

func ParseEvent(raw []byte) (slackevents.EventsAPIEvent, error) {
	//why but why
	slackevents.EventsAPIInnerEventMapping = map[slackevents.EventsAPIType]interface{}{
		slackevents.AppMention:      slack.MessageEvent{},
		slackevents.ReactionAdded:   slackevents.ReactionAddedEvent{},
		slackevents.URLVerification: slackevents.EventsAPIURLVerificationEvent{},
	}

	return slackevents.ParseEvent(json.RawMessage(raw), slackevents.OptionNoVerifyToken())
}

func (c *Client) VerifySignature(header http.Header, raw []byte) error {
	sv, err := slack.NewSecretsVerifier(header, c.signingSecret)
	if err != nil {
		return err
	}

	sv.Write(raw)
	sv.WithDebug(&DebugLogging{})
	return sv.Ensure()
}
