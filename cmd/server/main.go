package main

import (
	"github.com/EmmaO/UrlShortener/pkg/storage/firestore"
	"github.com/EmmaO/UrlShortener/pkg/getting"
	"log"
	"net/http"
	"github.com/EmmaO/UrlShortener/pkg/http/rest"
	"github.com/EmmaO/UrlShortener/pkg/adding"
)

func main() {
	repository, err := firestore.NewStorage()

	if err != nil {
		log.Fatal(err)
	}

	adder := adding.NewService(repository)
	getter := getting.NewService(repository)

	router := rest.Router(adder, getter)

	log.Fatal(http.ListenAndServe(":8080", router))
}
