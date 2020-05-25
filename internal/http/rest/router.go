package rest

import (
	"net/http"
	"encoding/json"

	"github.com/EmmaO/UrlShortener/internal/getting"
	"github.com/EmmaO/UrlShortener/internal/adding"

	"github.com/gorilla/mux"
)

//Router returns a mux router with endpoints registered against Handlers
func Router(addingService adding.Service, gettingService getting.Service) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/r/{hashKey}", redirectByHashKey(gettingService)).Methods("GET")
	router.HandleFunc("/{hashKey}", createNewShortenedURL(addingService)).Methods("POST")
	router.HandleFunc("/", createNewShortenedURLRandomHash(addingService)).Methods("POST")

	return router
}

func createNewShortenedURL(service adding.Service) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var hashKey = mux.Vars(r)["hashKey"]

		var requestBody ShortenedURLApiRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		documentID, err := service.AddShortenedURL(adding.ShortenedURLRequest{
			FullURL: requestBody.FullURL,
			HashKey: hashKey,
		})

		if err != nil {
			switch err {
			case adding.ErrAlreadyExists:
				w.WriteHeader(http.StatusConflict)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(documentID))
	}
}

func createNewShortenedURLRandomHash(service adding.Service) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody ShortenedURLApiRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		documentID, err := service.AddShortenedURLWithRandomHash(requestBody.FullURL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(documentID))
	}
}

func redirectByHashKey(service getting.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hashKey := mux.Vars(r)["hashKey"]

		var shortenedURL getting.ShortenedURL
		shortenedURL, err := service.GetShortenedURL(hashKey)

		if err != nil {
			if err == getting.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		http.Redirect(w, r, shortenedURL.FullURL, http.StatusMovedPermanently)
	}
}
