package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/behind-the-scenes/pkg/api"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
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
