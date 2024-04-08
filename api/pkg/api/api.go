package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rajatjindal/behind-the-scenes/api/pkg/posts"
	"github.com/rajatjindal/behind-the-scenes/api/pkg/webhook"
)

// Server is api server
type Server struct {
	Router *mux.Router
}

// New returns new server
func New() (*Server, error) {
	router := mux.NewRouter().StrictSlash(true)
	server := &Server{
		Router: router,
	}

	err := server.addRoutes()
	if err != nil {
		return nil, err
	}

	return server, nil
}

const uuidRegex = "[a-fA-F0-9]{8}-?[a-fA-F0-9]{4}-?4[a-fA-F0-9]{3}-?[8|9|aA|bB][a-fA-F0-9]{3}-?[a-fA-F0-9]{12}"

func (s *Server) addRoutes() error {
	s.Router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	eventhandler, err := webhook.NewHandler()
	if err != nil {
		return err
	}

	s.Router.Methods(http.MethodGet).Path("/api/runs-on").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runsOn, err := variables.Get("runs_on")
		if err != nil {
			logrus.Errorf("error marshalling %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write([]byte(runsOn))
	})

	s.Router.Methods(http.MethodGet).Path("/api/posts").HandlerFunc(posts.GetPostsHandler)
	s.Router.Methods(http.MethodGet).Path("/api/post/{postId}").HandlerFunc(posts.GetPostHandler)
	s.Router.Methods(http.MethodPost).Path("/api/post/{postId}/grapes").HandlerFunc(posts.IncrementGrapesHandler)
	s.Router.Methods(http.MethodPost).Path("/api/post/{postId}/hearts").HandlerFunc(posts.IncrementHeartsHandler)
	s.Router.Methods(http.MethodGet).Path("/api/post/{postId}/image/{imageId:" + uuidRegex + "}").HandlerFunc(posts.GetImageHandler)
	s.Router.Methods(http.MethodPost).Path("/api/slack").HandlerFunc(eventhandler.Handle)

	//add middleware after all endpoints are added
	s.Router.Use(CorsHandler())

	return nil
}
