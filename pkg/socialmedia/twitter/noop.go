package twitter

import (
	"context"
	"fmt"

	"github.com/rajatjindal/behind-the-scenes/pkg/creds"
	"github.com/rajatjindal/behind-the-scenes/pkg/socialmedia"
)

type NoopClient struct{}

// NewClient returns new twitter client
func NewNoopClient(credsProvider creds.Provider) (*NoopClient, error) {
	return &NoopClient{}, nil
}

var (
	_ socialmedia.Provider = &NoopClient{}
)

func (c *NoopClient) CreatePost(ctx context.Context, text string, images ...socialmedia.Image) error {
	fmt.Println(ctx, text, images)
	return nil
}

func (c *NoopClient) Name() string {
	return "noop"
}
