package webhook

import (
	"bytes"
	"context"

	"github.com/fermyon/spin/sdk/go/config"
	kv "github.com/fermyon/spin/sdk/go/key_value"
	"github.com/rajatjindal/pets-of-fermyon/pkg/bluesky"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	TriggerEmojiCodeKey = "trigger_on_emoji_code"
	MaxImagesInPost     = 4
)

func (s *Service) handleCallbackEvent(ctx context.Context, outerEvent slackevents.EventsAPIEvent) error {
	reactionAddedEvent, ok := outerEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
	if !ok {
		return nil
	}

	return s.handleReactionAddedEvent(ctx, reactionAddedEvent)
}

func (s *Service) handleReactionAddedEvent(ctx context.Context, event *slackevents.ReactionAddedEvent) error {
	ok, err := s.already_processed(event.Item.Timestamp)
	if err != nil {
		return err
	}

	// already processed
	if ok {
		return nil
	}

	triggerEmoji, err := config.Get(TriggerEmojiCodeKey)
	if err != nil {
		return err
	}

	if event.Reaction != triggerEmoji {
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

func (s *Service) already_processed(key string) (bool, error) {
	store, err := kv.Open("default")
	if err != nil {
		return true, err
	}
	defer kv.Close(store)

	return kv.Exists(store, key)
}

func (s *Service) set_already_processed(key string) error {
	store, err := kv.Open("default")
	if err != nil {
		return err
	}
	defer kv.Close(store)

	return kv.Set(store, key, []byte("true"))
}
