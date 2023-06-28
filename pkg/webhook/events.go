package webhook

import (
	"bytes"
	"context"
	"fmt"
	"time"

	kv "github.com/fermyon/spin/sdk/go/key_value"
	"github.com/rajatjindal/pets-of-fermyon/pkg/bluesky"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	TriggerEmojiCodeKey = "trigger_on_emoji_code"
	MaxImagesInPost     = 4
)

func (s *Handler) handleCallbackEvent(ctx context.Context, outerEvent slackevents.EventsAPIEvent) error {
	reactionAddedEvent, ok := outerEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
	if !ok {
		return nil
	}

	return s.handleReactionAddedEvent(ctx, reactionAddedEvent)
}

func (s *Handler) handleReactionAddedEvent(ctx context.Context, event *slackevents.ReactionAddedEvent) error {
	ok, err := s.already_processed(event.Item.Timestamp)
	if err != nil {
		return err
	}

	// already processed
	if ok {
		return nil
	}

	// IMP: safegaurd to allow app in specific channels only
	if event.Item.Channel != s.allowedChannelCode {
		return nil
	}

	if event.Reaction != s.triggerEmojiCode {
		return nil
	}

	resp, err := s.slack.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID:          event.Item.Channel,
		Latest:             event.Item.Timestamp,
		Limit:              1,
		Inclusive:          true,
		IncludeAllMetadata: true,
	})
	if err != nil {
		return nil
	}

	images := []string{}
	for _, msg := range resp.Messages {
		for _, file := range msg.Files {
			images = append(images, file.URLPrivateDownload)
		}
	}

	fmt.Printf("%s no of images %d\n", time.Now().Format(time.RFC3339), len(images))
	imagesToEmbed := []bluesky.Image{}
	for _, image := range images {
		var imageFile bytes.Buffer
		err = s.slack.GetFile(image, &imageFile)
		if err != nil {
			return err
		}

		imagesToEmbed = append(imagesToEmbed, imageFile.Bytes())

		if len(imagesToEmbed) == MaxImagesInPost {
			break
		}
	}

	err = s.bluesky.CreatePost(ctx, imagesToEmbed)
	if err != nil {
		return err
	}

	return nil
}

func (s *Handler) already_processed(key string) (bool, error) {
	store, err := kv.Open("default")
	if err != nil {
		return true, err
	}
	defer kv.Close(store)

	return kv.Exists(store, key)
}

func (s *Handler) set_already_processed(key string) error {
	store, err := kv.Open("default")
	if err != nil {
		return err
	}
	defer kv.Close(store)

	return kv.Set(store, key, []byte("true"))
}
