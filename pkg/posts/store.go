package posts

import (
	"encoding/json"
	"fmt"
	"strings"

	kv "github.com/fermyon/spin/sdk/go/key_value"
	"github.com/google/uuid"
	"github.com/rajatjindal/behind-the-scenes/pkg/logrus"
)

func Exists(id string) (bool, error) {
	store, err := kv.Open("default")
	if err != nil {
		return true, err
	}
	defer kv.Close(store)

	return kv.Exists(store, fmt.Sprintf("post:%s", id))
}

func StoreImageMeta(imageURL string) (string, error) {
	store, err := kv.Open("default")
	if err != nil {
		return "", err
	}
	defer kv.Close(store)

	imageId := uuid.New().String()
	skey := fmt.Sprintf("image:%s", imageId)
	logrus.Infof("adding key %s into store", skey)
	err = kv.Set(store, skey, []byte(imageURL))
	if err != nil {
		logrus.Infof("error when adding key %s into store %v", skey, err)
	}

	logrus.Infof("after adding key %s into store", skey)
	return imageId, err
}

func StorePost(post *Post) error {
	store, err := kv.Open("default")
	if err != nil {
		return err
	}
	defer kv.Close(store)

	raw, err := json.Marshal(post)
	if err != nil {
		return err
	}

	skey := fmt.Sprintf("post:%s", post.Timestamp)
	logrus.Infof("adding key %s into store", skey)
	err = kv.Set(store, skey, raw)
	if err != nil {
		logrus.Infof("error when adding key %s into store %v", skey, err)
	}
	logrus.Infof("after adding key %s into store", skey)
	return err
}

func DeleteAllPosts() error {
	store, err := kv.Open("default")
	if err != nil {
		return err
	}

	allKeys, err := kv.GetKeys(store)
	if err != nil {
		return err
	}

	for _, key := range allKeys {
		fmt.Printf("key is %s\n", key)
		if !strings.HasPrefix(key, "post:") {
			continue
		}

		fmt.Printf("deleting post with key %s\n", key)
		err = kv.Delete(store, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetAllPostsKeys() ([]string, error) {
	store, err := kv.Open("default")
	if err != nil {
		return nil, err
	}

	allKeys, err := kv.GetKeys(store)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, key := range allKeys {
		fmt.Printf("key is %s\n", key)
		if !strings.HasPrefix(key, "post:") {
			continue
		}

		keys = append(keys, strings.TrimPrefix(key, "post:"))
	}

	return keys, nil
}

func GetPostAsBytes(id string) ([]byte, error) {
	store, err := kv.Open("default")
	if err != nil {
		return nil, err
	}

	return kv.Get(store, fmt.Sprintf("post:%s", id))
}

func GetPost(id string) (*Post, error) {
	store, err := kv.Open("default")
	if err != nil {
		return nil, err
	}

	raw, err := kv.Get(store, fmt.Sprintf("post:%s", id))
	if err != nil {
		return nil, err
	}

	var post Post
	err = json.Unmarshal(raw, &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
