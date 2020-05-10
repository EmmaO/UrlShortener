package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//ShortenedURL contains details for redirecting from hash key to full URL
type ShortenedURL struct {
	HashKey       string
	FullURL       string
	ExpirationUTC time.Time
}

var firestoreClient firestore.Client

func main() {
	configureFirebaseClient()

	router := mux.NewRouter()
	router.HandleFunc("/{hashKey}", redirectByHashKey).Methods("GET")
}

func configureFirebaseClient() {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./saCreds.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()
}

func redirectByHashKey(w http.ResponseWriter, r *http.Request) {
	var hashKey = mux.Vars(r)["hashKey"]

	app, _ := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(""))
	firestoreClient, _ := app.Firestore(context.Background())

	docRef := firestoreClient.Doc(fmt.Sprintf("shorturls/%v", hashKey))
	doc, err := docRef.Get(r.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	var shortenedURL ShortenedURL
	if err := doc.DataTo(&shortenedURL); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	if shortenedURL.ExpirationUTC.Before(time.Now().UTC()) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, shortenedURL.FullURL, http.StatusMovedPermanently)
}
