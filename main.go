package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fermyon/spin/sdk/go/config"
	spinhttp "github.com/fermyon/spin/sdk/go/http"
	kv "github.com/fermyon/spin/sdk/go/key_value"
	"github.com/rajatjindal/pets-of-fermyon/pkg/bluesky"
	"github.com/rajatjindal/pets-of-fermyon/pkg/creds/kvcreds"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		eventHandler(w, r)
	})
}

func main() {}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//why but why
	slackevents.EventsAPIInnerEventMapping = map[slackevents.EventsAPIType]interface{}{
		slackevents.ReactionAdded: slackevents.ReactionAddedEvent{},
	}

	fmt.Println(string(body))
	fmt.Println(1)
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(2)
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	fmt.Println(3)
	if eventsAPIEvent.Type != slackevents.CallbackEvent {
		fmt.Fprintln(w, "OK")
		return
	}

	fmt.Println(4)
	event, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
	if !ok {
		fmt.Fprintln(w, "OK")
		return
	}

	fmt.Println(5)
	ok, err = already_processed(event.Item.Timestamp)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(6)
	// already processed
	if ok {
		fmt.Fprintln(w, "OK")
		return
	}

	fmt.Println(7)
	triggerEmoji, err := config.Get("trigger_on_emoji_code")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(8)
	if event.Reaction != triggerEmoji {
		fmt.Fprintln(w, "OK")
		return
	}

	fmt.Println(9)
	client := spinhttp.NewClient()
	api := slack.New("xoxb-4430321508054-5503179271856-u2b6oMyPEhqnZunK0PBy1e32", slack.OptionHTTPClient(client))
	resp, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID:          event.Item.Channel,
		Latest:             event.Item.Timestamp,
		Limit:              1,
		Inclusive:          true,
		IncludeAllMetadata: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(10)
	image := ""
	for index, msg := range resp.Messages {
		fmt.Println(11)
		fmt.Printf("%#v", msg)
		fmt.Println()
		for _, file := range msg.Files {
			fmt.Printf("index: %d, attach: %#v", index, file.URLPrivateDownload)
			image = file.URLPrivateDownload
			fmt.Println()
		}
	}

	var imageFile bytes.Buffer
	err = api.GetFile(image, &imageFile)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	credsProvider := kvcreds.Provider()
	bluesky, err := bluesky.NewClient(client, credsProvider)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bluesky.CreatePost(ctx, imageFile.Bytes())
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func already_processed(key string) (bool, error) {
	store, err := kv.Open("default")
	if err != nil {
		return true, err
	}
	defer kv.Close(store)

	return kv.Exists(store, key)
}

func set_already_processed(key string) error {
	store, err := kv.Open("default")
	if err != nil {
		return err
	}
	defer kv.Close(store)

	return kv.Set(store, key, []byte("true"))
}
