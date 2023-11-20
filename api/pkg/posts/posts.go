package posts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/variables"
	"github.com/gorilla/mux"
	"github.com/rajatjindal/behind-the-scenes/api/pkg/slack"
	"github.com/sirupsen/logrus"
)

type Post struct {
	Msg       string            `json:"msg"`
	Timestamp string            `json:"timestamp"`
	ImageIds  []string          `json:"imageIds"`
	ImageMap  map[string]string `json:"imageMap,omitempty"`
	Approved  bool              `json:"approved"`
	Grapes    int               `json:"grapes"`
	Hearts    int               `json:"hearts"`
}

func IncrementGrapesHandler(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["postId"]
	post, err := GetPost(postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.Grapes = post.Grapes + 1

	err = StorePost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "OK")
}

func IncrementHeartsHandler(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["postId"]
	post, err := GetPost(postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.Hearts = post.Hearts + 1

	err = StorePost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "OK")
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

	apikey, err := variables.Get("bts_admin_api_key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !strings.Contains(r.Header.Get("authorization"), apikey) {
		logrus.Error("token mismatch when trying to delete all posts")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err = DeleteAllPosts()
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

	if !post.Approved {
		http.Error(w, "post not approved", http.StatusNotFound)
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

	client := spinhttp.NewClient()
	slackClient, err := slack.NewClient(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var imageFile bytes.Buffer
	err = slackClient.GetFile(imageURL, &imageFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(imageFile.Bytes())
}
