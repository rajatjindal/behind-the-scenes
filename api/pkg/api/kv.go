package api

import (
	"net/http"
	"strings"

	kv "github.com/fermyon/spin-go-sdk/kv"
	"github.com/fermyon/spin-go-sdk/variables"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func DeleteFromKV(w http.ResponseWriter, r *http.Request) {
	logrus.Info("delete from kv handler")

	apikey, err := variables.Get("bts_admin_api_key")
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !strings.Contains(r.Header.Get("authorization"), apikey) {
		logrus.Error("token mismatch when trying to delete all posts")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	key := mux.Vars(r)["key"]
	if key == "" {
		logrus.Error("key not provided to delete")
		http.Error(w, "key to delete not provided", http.StatusBadRequest)
		return
	}
	logrus.Infof("deleting %s from kv handler", key)

	store, err := kv.OpenStore("default")
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer store.Close()

	err = store.Delete(key)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
