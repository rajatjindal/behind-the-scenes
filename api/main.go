package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/rajatjindal/behind-the-scenes/api/pkg/api"
	"github.com/sirupsen/logrus"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("starting handle of %s", r.URL.Path)
		s, err := api.New()
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		s.Router.ServeHTTP(w, r)
	})
}

func main() {}
