package services

import (
	"time"
	"url-shortener-go/internal/repository"
)

type URLService interface {
	CreateShortURL(long string) (string, error)
	GetLongURL(short string) (string, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) CreateShortURL(long string) (string, error) {
	short := generateShortID()
	err := s.repo.StoreURL(short, long)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (s *urlService) GetLongURL(short string) (string, error) {
	return s.repo.GetLongURL(short)
}

// generate a random short ID string.
func generateShortID() string {
	// For example, generate a random hex string of length 5
	const letters = "0123456789abcdef"
	b := make([]byte, 5)
	for i := range b {
		b[i] = letters[randInt(len(letters))]
	}
	return string(b)
}

// return a random int between 0 and max-1
func randInt(max int) int {
	return int(time.Now().UnixNano() % int64(max))
}
