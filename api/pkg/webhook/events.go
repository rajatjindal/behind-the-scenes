package webhook

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	kv "github.com/fermyon/spin-go-sdk/kv"
	"github.com/google/uuid"
	"github.com/rajatjindal/behind-the-scenes/api/pkg/posts"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	TriggerEmojiCodeKey = "trigger_on_emoji_code"
	self                = "U05HDSRLEF7"
)

func (s *Handler) getSelfUser() (string, error) {
	store, err := kv.OpenStore("default")
	if err != nil {
		return "", err
	}

	val, err := store.Get("self-user-id")
	if err != nil {
		return "", err
	}

	if val != nil {
		return string(val), nil
	}

	resp, err := s.slack.GetUserIdentity()
	if err != nil {
		return "", err
	}

	_ = store.Set("self-user-id", []byte(resp.User.ID))
	return resp.User.ID, nil
}

func (s *Handler) handleCallbackEvent(ctx context.Context, outerEvent slackevents.EventsAPIEvent) error {
	logrus.Info("starting handleCallbackEvent")
	appMentionEvent, ok := outerEvent.InnerEvent.Data.(*slack.MessageEvent)
	if ok {
		return s.handleAppMentionEvent(ctx, appMentionEvent)
	}

	reactionAddedEvent, ok := outerEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
	if ok {
		return s.handleReactionAddedEvent(ctx, reactionAddedEvent)
	}

	return fmt.Errorf("unsupported event")
}

func (s *Handler) handleAppMentionEvent(ctx context.Context, event *slack.MessageEvent) error {
	logrus.Info("starting handleAppMentionEvent")
	if exists, err := posts.Exists(event.Timestamp); exists || err != nil {
		return nil
	}

	// IMP: safegaurd to allow app in specific channels only
	if event.Channel != s.allowedChannelCode {
		logrus.Infof("channel %q is not allowed %q", event.Channel, s.allowedChannelCode)
		return nil
	}

	if !verifySignoffFromEvent(event) {
		return nil
	}

	imageIdsMap := map[string]string{}
	imageIds := []string{}
	for _, file := range event.Files {
		if file.Filetype == "mp4" {
			continue
		}

		imageId := uuid.New().String()

		imageIdsMap[imageId] = file.URLPrivateDownload
		imageIds = append(imageIds, imageId)
	}

	logrus.Infof("images to add: %v", imageIds)
	if len(imageIds) == 0 {
		// no images to post
		return nil
	}

	post := &posts.Post{
		Msg:       event.Text,
		ImageIds:  imageIds,
		ImageMap:  imageIdsMap,
		Timestamp: event.Timestamp,
		Approved:  verifySignoffFromEvent(event),
	}

	err := posts.StorePost(post)
	if err != nil {
		logrus.Errorf("failed to store post. err: %v", err)
	}

	return nil
}

func (s *Handler) handleReactionAddedEvent(ctx context.Context, event *slackevents.ReactionAddedEvent) error {
	logrus.Info("starting handleReactionAddedEvent")
	// IMP: safegaurd to allow app in specific channels only
	if event.Item.Channel != s.allowedChannelCode {
		logrus.Infof("channel %q is not allowed. Allowed channel is %q", event.Item.Channel, s.allowedChannelCode)
		return nil
	}

	resp, err := s.slack.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID:          event.Item.Channel,
		Latest:             event.Item.Timestamp,
		Limit:              1,
		Inclusive:          true,
		IncludeAllMetadata: true,
	})
	if err != nil || len(resp.Messages) == 0 {
		return nil
	}
	logrus.Info("able to get the conversation history successfully")

	// check if all users who were tagged in original msg have
	// added the signoff emoji
	msg := resp.Messages[0]
	signedOff := verifySignoffFromMessage(msg)
	if !signedOff {
		return nil
	}

	//signoff is done
	msgTxt := msg.Text
	imageURLs := []string{}

	for _, file := range msg.Files {
		if file.Filetype == "mp4" {
			continue
		}

		imageURLs = append(imageURLs, file.URLPrivateDownload)
	}

	logrus.Infof("number of images in event %d", len(imageURLs))
	if len(imageURLs) == 0 {
		return nil
	}

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
		Approved:  verifySignoffFromMessage(msg),
	}

	err = posts.StorePost(post)
	if err != nil {
		logrus.Errorf("failed to set_already_processed. err: %v", err)
	}

	return nil
}

func getTaggedUsers(text string) []string {
	r := regexp.MustCompile(`<@U\w+>`)
	matches := r.FindAllString(text, -1)

	users := []string{}
	userMap := map[string]struct{}{}
	for _, m := range matches {
		m = strings.TrimLeft(m, "<@")
		m = strings.TrimRight(m, ">")

		if _, exists := userMap[m]; exists {
			continue
		}

		userMap[m] = struct{}{}
		users = append(users, m)
	}

	return users
}

func verifySignoffFromEvent(event *slack.MessageEvent) bool {
	if !strings.Contains(event.Text, "/signoff") {
		return false
	}

	taggedUsers := getTaggedUsers(event.Text)
	approvedUsers := []string{event.User, self} // assume already approved by author and bts bot
	for _, reaction := range event.Reactions {
		if reaction.Name != "squirrel" {
			continue
		}

		approvedUsers = append(approvedUsers, reaction.Users...)
	}

	return isApproved(taggedUsers, approvedUsers)
}

func verifySignoffFromMessage(msg slack.Message) bool {
	if !strings.Contains(msg.Text, "/signoff") {
		return false
	}

	taggedUsers := getTaggedUsers(msg.Text)
	approvedUsers := []string{msg.User, self} // assume already approved by author
	for _, reaction := range msg.Reactions {
		if reaction.Name != "squirrel" {
			continue
		}

		approvedUsers = append(approvedUsers, reaction.Users...)
	}

	return isApproved(taggedUsers, approvedUsers)
}

func isApproved(required []string, approvedUsers []string) bool {
	logrus.Infof("required: %v\napproved: %v", required, approvedUsers)
	for _, item := range required {
		if !contains(approvedUsers, item) {
			return false
		}
	}

	return true
}

func contains(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}

	return false
}
