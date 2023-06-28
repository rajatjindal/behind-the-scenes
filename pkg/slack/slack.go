package slack

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Client struct {
	*slack.Client
}

func New(httpclient *http.Client, token string) *Client {
	client := slack.New(token, slack.OptionHTTPClient(httpclient))
	return &Client{
		Client: client,
	}
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
