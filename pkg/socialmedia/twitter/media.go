package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rajatjindal/behind-the-scenes/pkg/logrus"
)

type UploadMediaResp struct {
	MediaId int64 `json:"media_id"`
}

func (t *RealClient) upload(ctx context.Context, image []byte) (int64, error) {
	logrus.Info("entered upload of image func")
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("media", "image.jpeg")
	if err != nil {
		return 0, err
	}

	_, err = io.Copy(part, bytes.NewReader(image))
	if err != nil {
		return 0, err
	}
	writer.Close()

	logrus.Info("creating request to upload image")
	req, err := http.NewRequest(http.MethodPost, "https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image", &body)
	if err != nil {
		return 0, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	logrus.Info("starting the upload of image")
	resp, err := t.oauth.Do(req)
	if err != nil {
		return 0, err
	}

	logrus.Info("reading response body")
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	logrus.Info("checking status code")
	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("expected status code: %d, got: %d. body: %s", http.StatusOK, resp.StatusCode, string(raw))
		return 0, fmt.Errorf("expected status code: %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	logrus.Info("parsing response")
	media := UploadMediaResp{}

	err = json.Unmarshal(raw, &media)
	if err != nil {
		return 0, err
	}

	logrus.Info("returning media id")
	return media.MediaId, nil
}
