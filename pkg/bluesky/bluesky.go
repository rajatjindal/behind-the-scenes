package bluesky

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/rajatjindal/pets-of-fermyon/pkg/creds"
)

const (
	base = "https://bsky.social"
)

type BlueSky struct {
	xrpcc *xrpc.Client
}

func NewClient(client *http.Client, credsProvider creds.Provider) (*BlueSky, error) {
	credentials, err := credsProvider.GetCredentials("bluesky")
	if err != nil {
		return nil, err
	}

	xrpcc := &xrpc.Client{
		Auth: &xrpc.AuthInfo{
			Handle: credentials["handle"],
		},
		Client: client,
		Host:   base,
	}

	auth, err := comatproto.ServerCreateSession(context.TODO(), xrpcc, &comatproto.ServerCreateSession_Input{
		Identifier: xrpcc.Auth.Handle,
		Password:   credentials["password"],
	})
	if err != nil {
		return nil, err
	}

	xrpcc.Auth.AccessJwt = auth.AccessJwt
	xrpcc.Auth.RefreshJwt = auth.RefreshJwt
	xrpcc.Auth.Did = auth.Did
	xrpcc.Auth.Handle = auth.Handle

	return &BlueSky{
		xrpcc: xrpcc,
	}, nil
}

func (b *BlueSky) CreatePost(ctx context.Context, imageFile []byte) error {
	post, err := b.format(ctx, imageFile)
	if err != nil {
		return err
	}

	_, err = comatproto.RepoCreateRecord(ctx, b.xrpcc, post)
	return err
}

func getImage(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image")
	}

	return io.ReadAll(resp.Body)
}

func (b *BlueSky) format(ctx context.Context, imageFile []byte) (*comatproto.RepoCreateRecord_Input, error) {
	postMsg := "somemsg @rajatjindal.bsky.social is testing stuff here"
	post := &appbsky.FeedPost{
		Text:      postMsg,
		CreatedAt: time.Now().Format(time.RFC3339),
		Facets:    DetectFacets(postMsg),
	}

	og, err := b.getEmbedData(ctx, imageFile)
	if err != nil {
		return nil, err
	}

	//add embed image info
	post.Embed = og

	return &comatproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       b.xrpcc.Auth.Did,
		Record: &util.LexiconTypeDecoder{
			Val: post,
		},
	}, nil
}

func (b *BlueSky) Name() string {
	return "bluesky"
}

func (b *BlueSky) getEmbedData(ctx context.Context, imageFile []byte) (*appbsky.FeedPost_Embed, error) {
	blob, err := comatproto.RepoUploadBlob(ctx, b.xrpcc, bytes.NewReader(imageFile))
	if err != nil {
		return nil, err
	}

	return &appbsky.FeedPost_Embed{
		EmbedImages: &appbsky.EmbedImages{
			LexiconTypeID: "",
			Images: []*appbsky.EmbedImages_Image{
				{
					Alt: "",
					Image: &util.LexBlob{
						Ref:      blob.Blob.Ref,
						MimeType: "image/jpeg",
					},
				},
			},
		},
	}, nil
}
