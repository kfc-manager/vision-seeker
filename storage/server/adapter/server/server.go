package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kfc-manager/vision-seeker/storage/server/service/data"
)

type Server interface {
	Listen() error
}

type server struct {
	router *http.ServeMux
	port   int
	token  string
	data   data.Service
}

func New(port int, token string, service data.Service) *server {
	s := &server{
		port:  port,
		token: token,
		data:  service,
	}
	router := http.NewServeMux()
	router.HandleFunc("/dataset/{name}", s.handleDataset)
	router.HandleFunc("/{hash}", s.handleImg)
	s.router = router
	return s
}

func (s *server) Listen() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}

func (s *server) handleAuth(w http.ResponseWriter, r *http.Request) error {
	if len(r.Header.Get("X-Custom-Token")) < 1 {
		err := errors.New("Unauthorized")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}
	if r.Header.Get("X-Custom-Token") != s.token {
		err := errors.New("Forbidden")
		http.Error(w, err.Error(), http.StatusForbidden)
		return err
	}
	return nil
}

func (s *server) handleDataset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Write([]byte{})
		return
	}

	if r.Method == http.MethodPost {
		err := s.handleAuth(w, r)
		if err != nil {
			return
		}

		name := r.PathValue("name")
		if len(name) < 1 {
			http.Error(w, "path value 'name' missing", http.StatusBadRequest)
			return
		}

		id, err := s.data.StoreDataset(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(id))
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (s *server) handleImg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Write([]byte{})
		return
	}

	if r.Method == http.MethodPost {
		err := s.handleAuth(w, r)
		if err != nil {
			return
		}

		id := r.Header.Get("X-Custom-Dataset")
		if len(id) < 1 {
			http.Error(w, "missing header 'X-Custom-Dataset'", http.StatusBadRequest)
			return
		}
		hash := r.PathValue("hash")
		if len(hash) != 64 {
			http.Error(w, "path value 'hash' invalid length", http.StatusBadRequest)
			return
		}
		label := r.Header.Get("X-Custom-Label")

		if r.Body == nil {
			http.Error(w, "missing request body", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "couldn't load request body", http.StatusInternalServerError)
			return
		}

		err = s.data.StoreImg(hash, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.data.StoreMetadata(id, hash, label)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
