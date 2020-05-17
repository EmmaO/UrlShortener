package firestore

import (
	"context"
	"time"

	"github.com/EmmaO/UrlShortener/pkg/adding"
	"github.com/EmmaO/UrlShortener/pkg/getting"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	//ShortURLCollection is the collection name for short urls
	ShortURLCollection = "shorturls"
)

//Storage stores shortened URLs in a firestore database
type Storage struct {
	db *firestore.Client
}

//NewStorage creates a new firetore storage
func NewStorage() (*Storage, error) {
	s := new(Storage)
	db, err := createFirestoreClient()
	s.db = db
	if err != nil {
		//TODO: return ErrRepositoryConfiguration instead
		return s, err
	}

	return s, nil
}

func createFirestoreClient() (*firestore.Client, error) {
	var firestoreClient *firestore.Client

	ctx := context.Background()

	// Use a service account
	sa := option.WithCredentialsFile("../../saCreds.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return firestoreClient, err
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		return firestoreClient, err
	}

	return firestoreClient, nil
}

//GetShortenedURL returns the ShortenedURL identified by the provided hashKey
func (storage Storage) GetShortenedURL(hashKey string) (getting.ShortenedURL, error) {
	doc, err := storage.db.Collection(ShortURLCollection).Doc(hashKey).Get(context.Background())
	var returnValue getting.ShortenedURL

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return returnValue, getting.ErrNotFound
		}

		return returnValue, getting.ErrRepositoryUnknown
	}

	var shortenedURL ShortenedURL
	if err := doc.DataTo(&shortenedURL); err != nil {
		return returnValue, getting.ErrRepositoryUnknown
	}

	returnValue.HashKey = shortenedURL.HashKey
	returnValue.FullURL = shortenedURL.FullURL
	returnValue.ExpirationUTC = shortenedURL.ExpirationUTC

	return returnValue, nil
}

//AddShortenedURL adds a new shortenedUrl to the repository
func (storage Storage) AddShortenedURL(shortenedURL adding.ShortenedURL) error {
	//expect not found error
	existing, err := storage.GetShortenedURL(shortenedURL.HashKey)
	if err == nil {
		if existing.ExpirationUTC.After(time.Now().UTC()) {
			// return already exists in link is present and not expired
			return adding.ErrAlreadyExists
		}
	} else if err != getting.ErrNotFound {
		//something else went wrong
		return adding.ErrRepositoryUnknown
	}

	var repositoryModel ShortenedURL
	repositoryModel.HashKey = shortenedURL.HashKey
	repositoryModel.FullURL = shortenedURL.FullURL
	repositoryModel.ExpirationUTC = shortenedURL.ExpirationUTC

	_, err = storage.db.Collection(ShortURLCollection).Doc(repositoryModel.HashKey).Set(context.Background(), repositoryModel)

	if err != nil {
		return err
	}

	return nil
}
