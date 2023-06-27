package slack

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rajatjindal/behind-the-scenes/pkg/creds"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Client struct {
	signingSecret string
	*slack.Client
}

func NewClient(httpclient *http.Client, credsProvider creds.Provider) (*Client, error) {
	creds, err := credsProvider.GetCredentials("slack")
	if err != nil {
		return nil, err
	}

	token := creds["token"]
	if token == "" {
		return nil, fmt.Errorf("slack token not found")
	}

	signingSecret := creds["signingSecret"]
	if signingSecret == "" {
		return nil, fmt.Errorf("slack signing secret not found")
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
		slackevents.ReactionAdded:   slackevents.ReactionAddedEvent{},
		slackevents.URLVerification: slackevents.EventsAPIURLVerificationEvent{},
	}

	return slackevents.ParseEvent(json.RawMessage(raw), slackevents.OptionNoVerifyToken())
}

func URLVerificationHandler(w http.ResponseWriter, raw []byte) {
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

func (c *Client) VerifySignature(header http.Header, raw []byte) error {
	sv, err := slack.NewSecretsVerifier(header, c.signingSecret)
	if err != nil {
		return err
	}

	sv.Write(raw)
	sv.WithDebug(&DebugLogging{})
	return sv.Ensure()
}
