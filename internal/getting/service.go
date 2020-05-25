package getting

import "time"

//Service used for getting shortened URLs from the repository
type Service interface {
	GetShortenedURL(hashKey string) (ShortenedURL, error)
}

//Repository with methods for getting ShortenedURLs
type Repository interface {
	GetShortenedURL(haskKey string) (ShortenedURL, error)
}

type service struct {
	repository Repository
}

//NewService returns a getting service
func NewService(r Repository) Service {
	return &service{r}
}

//GetShortenedURL returns a URL from the database if it is found and unexpired
func (service service) GetShortenedURL(hashKey string) (ShortenedURL, error) {
	var shortenedURL ShortenedURL
	shortenedURL, err := service.repository.GetShortenedURL(hashKey)

	if err != nil {
		return shortenedURL, err
	}

	if shortenedURL.ExpirationUTC.Before(time.Now().UTC()) {
		return shortenedURL, ErrNotFound
	}

	return shortenedURL, nil
}
