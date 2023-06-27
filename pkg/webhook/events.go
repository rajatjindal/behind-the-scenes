package webhook

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rajatjindal/behind-the-scenes/pkg/logrus"
	"github.com/rajatjindal/behind-the-scenes/pkg/posts"
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
		logrus.Info("not a reactionAddedEvent")
		return nil
	}

	return s.handleReactionAddedEvent(ctx, reactionAddedEvent)
}

func (s *Handler) handleReactionAddedEvent(ctx context.Context, event *slackevents.ReactionAddedEvent) error {
	ok, err := posts.Exists(event.Item.Timestamp)
	if err != nil {
		logrus.Info("error checking if already processed this msg")
		return err
	}

	// already processed
	if ok {
		logrus.Info("already processed this msg")
		return nil
	}

	// IMP: safegaurd to allow app in specific channels only
	if event.Item.Channel != s.allowedChannelCode {
		logrus.Infof("channel %q is not allowed %q", event.Item.Channel, s.allowedChannelCode)
		return nil
	}

	if event.Reaction != s.triggerEmojiCode {
		logrus.Info("emoji %q is not the triggerEmoji %q", event.Reaction, s.triggerEmojiCode)
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
	logrus.Info("able to get the conversation history successfully")

	msgTxt := ""
	imageURLs := []string{}
	for _, msg := range resp.Messages {
		msgTxt = msg.Text
		logrus.Infof("slack msg is %s", msg.Text)

		// if author added the /signoff, it means they are giving consent
		// for this image to be posted to social media. ignore otherwise.
		if !strings.Contains(msg.Text, "/signoff") {
			continue
		}

		for _, file := range msg.Files {
			imageURLs = append(imageURLs, file.URLPrivateDownload)
		}
	}

	logrus.Infof("number of images in event %d", len(imageURLs))
	if len(imageURLs) == 0 {
		return nil
	}

	// imagesToEmbed := []socialmedia.Image{}
	// for _, image := range images {
	// 	var imageFile bytes.Buffer
	// 	err = s.slack.GetFile(image, &imageFile)
	// 	if err != nil {
	// 		logrus.Infof("failed to get the file. error: %v", err)
	// 		return err
	// 	}
	//
	// 	imagesToEmbed = append(imagesToEmbed, imageFile.Bytes())
	//
	// 	if len(imagesToEmbed) == MaxImagesInPost {
	// 		break
	// 	}
	// }
	//
	// logrus.Info("posting to socialmedia")
	// for _, sm := range s.socials {
	// 	err = sm.CreatePost(ctx, "", imagesToEmbed...)
	// 	if err != nil {
	// 		logrus.Errorf("failed to upload to %s. err: %v", sm.Name(), err)
	// 	}
	// }

	imageIdsMap := map[string]string{}
	imageIds := []string{}
	for _, imageURL := range imageURLs {
		imageId := uuid.New().String()

		imageIdsMap[imageId] = imageURL
		imageIds = append(imageIds, imageId)
	}

	post := &posts.Post{
		Msg:       msgTxt,
		ImageIds:  imageIds,
		ImageMap:  imageIdsMap,
		Timestamp: event.Item.Timestamp,
	}

	err = posts.StorePost(post)
	if err != nil {
		logrus.Errorf("failed to set_already_processed. err: %v", err)
	}

	return nil
}
