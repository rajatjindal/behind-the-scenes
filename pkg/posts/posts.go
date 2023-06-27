package posts

import (
	"bytes"
	"encoding/json"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/gorilla/mux"
	"github.com/rajatjindal/behind-the-scenes/pkg/creds/kvcreds"
	"github.com/rajatjindal/behind-the-scenes/pkg/logrus"
	"github.com/rajatjindal/behind-the-scenes/pkg/slack"
)

type Post struct {
	Msg       string            `json:"msg"`
	Timestamp string            `json:"timestamp"`
	ImageIds  []string          `json:"imageIds"`
	ImageMap  map[string]string `json:"imageMap,omitempty"`
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("posts handler")
	keys, err := GetAllPostsKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := json.Marshal(keys)
	if err != nil {
		logrus.Errorf("error marshalling %s", raw)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(raw)
}

func DeleteAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("delete all posts handler")
	err := DeleteAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := GetPost(mux.Vars(r)["postId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//remove images map from the data to avoid exposing slack urls
	post.ImageMap = nil
	raw, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(raw)
}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["postId"]
	imageId := mux.Vars(r)["imageId"]
	if postId == "" || imageId == "" {
		http.Error(w, "both postId and imageId are required", http.StatusInternalServerError)
		return
	}

	post, err := GetPost(postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageURL := post.ImageMap[imageId]
	logrus.Infof("going to fetch image from %s", imageURL)

	credsProvider := kvcreds.Provider()
	client := spinhttp.NewClient()

	slackClient, err := slack.NewClient(client, credsProvider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var imageFile bytes.Buffer
	err = slackClient.GetFile(imageURL, &imageFile)
	w.Write(imageFile.Bytes())
}
