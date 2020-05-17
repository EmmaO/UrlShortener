package adding

import (
	"github.com/EmmaO/UrlShortener/pkg/stringextensions"
	"time"
)

const (
	//DefaultExpirationDays notes the number of days to wait before marking a link as expired
	DefaultExpirationDays = 7

	//MaxHashGenerationRetryCount marks the maximum number of attempts the application will make to generate a new
	//unique hash before erroring
	MaxHashGenerationRetryCount = 10
)

//Service used for adding shortened URLs to the repository
type Service interface {
	AddShortenedURL(shortenedURL ShortenedURLRequest) (string, error)
	AddShortenedURLWithRandomHash(fullURL string) (string, error)
}

//Repository with methods for adding ShortenedURLs
type Repository interface {
	//AddShortenedURL adds shortened URLs to the repository
	AddShortenedURL(shortenedURL ShortenedURL) error
}

type service struct {
	repository Repository
}

//NewService returns an adding service
func NewService(r Repository) Service {
	return &service{r}
}

//AddShortenedURL adds a shortened URL with a default expiry time
func (service service) AddShortenedURL(request ShortenedURLRequest) (string, error) {
	//expect not found error
	shortenedURL := ShortenedURL{
		FullURL:       request.FullURL,
		HashKey:       request.HashKey,
		ExpirationUTC: time.Now().AddDate(0, 0, DefaultExpirationDays),
	}

	err := service.repository.AddShortenedURL(shortenedURL)

	if err != nil {
		return "", err
	}

	return shortenedURL.HashKey, nil
}

//AddShortenedURLWithRandomHash generates a random hash pairs it with the submitted URL and a default
//expiry time to create a shortened URL. This is then added to the repository
func (service service) AddShortenedURLWithRandomHash(fullURL string) (string, error) {
	shortenedURL := ShortenedURLRequest{
		FullURL: fullURL,
	}

	retryCount := 0
	for {
		shortenedURL.HashKey = stringextensions.RandomAlphaNumeric(7)
		_, err := service.AddShortenedURL(shortenedURL)

		if err != nil {
			if err == ErrAlreadyExists && retryCount > MaxHashGenerationRetryCount {
				return "", ErrHashGenerationFailed
			} else if err == ErrAlreadyExists {
				retryCount++
				continue
			} else {
				return "", err
			}
		}

		break
	}

	return shortenedURL.HashKey, nil
}
