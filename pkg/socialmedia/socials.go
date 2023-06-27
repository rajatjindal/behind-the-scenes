package socialmedia

import "context"

type Image []byte

type Provider interface {
	Name() string
	CreatePost(ctx context.Context, text string, images ...Image) error
}
