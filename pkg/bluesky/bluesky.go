package bluesky

import (
	"bytes"
	"context"
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

type Image []byte

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

func (b *BlueSky) Name() string {
	return "bluesky"
}

func (b *BlueSky) CreatePost(ctx context.Context, images []Image) error {
	post, err := b.format(ctx, images)
	if err != nil {
		return err
	}

	_, err = comatproto.RepoCreateRecord(ctx, b.xrpcc, post)
	return err
}

func (b *BlueSky) format(ctx context.Context, images []Image) (*comatproto.RepoCreateRecord_Input, error) {
	postMsg := ""
	post := &appbsky.FeedPost{
		Text:      postMsg,
		CreatedAt: time.Now().Format(time.RFC3339),
		Facets:    DetectFacets(postMsg),
	}

	og, err := b.getEmbedData(ctx, images)
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

func (b *BlueSky) getEmbedData(ctx context.Context, images []Image) (*appbsky.FeedPost_Embed, error) {
	embeddedImages := []*appbsky.EmbedImages_Image{}

	for _, image := range images {
		blob, err := comatproto.RepoUploadBlob(ctx, b.xrpcc, bytes.NewReader(image))
		if err != nil {
			return nil, err
		}

		embeddedImages = append(embeddedImages, &appbsky.EmbedImages_Image{
			Alt: "",
			Image: &util.LexBlob{
				Ref:      blob.Blob.Ref,
				MimeType: "image/jpeg",
			},
		})
	}

	return &appbsky.FeedPost_Embed{
		EmbedImages: &appbsky.EmbedImages{
			LexiconTypeID: "",
			Images:        embeddedImages,
		},
	}, nil
}
