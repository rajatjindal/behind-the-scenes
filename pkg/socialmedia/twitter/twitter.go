package twitter

import (
	"context"
	"errors"
	"net/http"

	//lint:ignore SA1019 ignore this for now!
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/rajatjindal/behind-the-scenes/pkg/creds"
	"github.com/rajatjindal/behind-the-scenes/pkg/logrus"
	"github.com/rajatjindal/behind-the-scenes/pkg/socialmedia"
)

// ErrTweetFailed is for failed tweet
var ErrTweetFailed = errors.New("failed to send tweet")

type RealClient struct {
	twitter *twitter.Client
	oauth   *http.Client
}

var (
	_ socialmedia.Provider = &RealClient{}
)

// NewClient returns new twitter client
func NewClient(client *http.Client, credsProvider creds.Provider) (*RealClient, error) {
	credentials, err := credsProvider.GetCredentials("twitter")
	if err != nil {
		return nil, err
	}

	config := oauth1.NewConfig(credentials["consumerKey"], credentials["consumerToken"])
	token := oauth1.NewToken(credentials["token"], credentials["tokenSecret"])
	httpClient := config.Client(context.WithValue(oauth1.NoContext, oauth1.HTTPClient, client), token)

	return &RealClient{
		twitter: twitter.NewClient(httpClient),
		oauth:   httpClient,
	}, nil
}

func (c *RealClient) CreatePost(ctx context.Context, prefix string, images ...socialmedia.Image) error {
	mediaIds := []int64{}
	for _, image := range images {
		id, err := c.upload(ctx, image)
		if err != nil {
			return err
		}

		mediaIds = append(mediaIds, id)
	}

	params := twitter.StatusUpdateParams{
		MediaIds: mediaIds,
	}

	_, _, err := c.twitter.Statuses.Update("", &params)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (c *RealClient) Name() string {
	return "twitter"
}
