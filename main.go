package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
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

var firestoreClient *firestore.Client

func main() {
	configureFirebaseClient()

	router := mux.NewRouter()
	router.HandleFunc("/r/{hashKey}", redirectByHashKey).Methods("GET")
	router.HandleFunc("/{hashKey}", createNewShortenedURL).Methods("POST")
	router.HandleFunc("/", createNewShortenedURLRandomHash).Methods("POST")

	http.ListenAndServe(":8080", router)
}

func configureFirebaseClient() {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./saCreds.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func redirectByHashKey(w http.ResponseWriter, r *http.Request) {
	hashKey := mux.Vars(r)["hashKey"]

	doc, err := firestoreClient.Collection("shorturls").Doc(hashKey).Get(context.Background())
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

type CreateShortenedUrlRequest struct {
	FullURL string
}

func createNewShortenedURLRandomHash(w http.ResponseWriter, r *http.Request) {

	var requestBody CreateShortenedUrlRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.FullURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := 0; i < 100; i++ {
		hashKey := generateRandomAlphaNumericString()

		documentID, err := tryAddShortenedURL(hashKey, requestBody.FullURL)
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(documentID))
			return
		} else if status.Code(err) != codes.AlreadyExists {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Unable to find unique hash after 100 attempts. Crazy."))
}

func generateRandomAlphaNumericString() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	returnValue := ""
	for i := 0; i < 7; i++ {
		returnValue += string(chars[rand.Intn(len(chars))])
	}

	return returnValue
}

func createNewShortenedURL(w http.ResponseWriter, r *http.Request) {
	var hashKey = mux.Vars(r)["hashKey"]

	var requestBody CreateShortenedUrlRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.FullURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	documentId, err := tryAddShortenedURL(hashKey, requestBody.FullURL)
	if err != nil {
		switch status.Code(err) {
		case codes.AlreadyExists:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(documentId))
}

func tryAddShortenedURL(hashKey string, fullURL string) (string, error) {
	//expect not found error
	ctx := context.Background()
	_, err := firestoreClient.Collection("shorturls").Doc(hashKey).Get(ctx)
	if err == nil {
		//found existing record with that hashKey
		return "", status.Error(codes.AlreadyExists, "key in use")
	} else if status.Code(err) != codes.NotFound {
		//something else went wrong
		return "", err
	}

	_, err = firestoreClient.Collection("shorturls").Doc(hashKey).Set(ctx, ShortenedURL{
		HashKey:       hashKey,
		FullURL:       fullURL,
		ExpirationUTC: time.Now().AddDate(0, 0, 7),
	})

	if err != nil {
		return "", err
	}

	return hashKey, nil
}
