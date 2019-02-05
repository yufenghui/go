package oauth

import (
	"errors"
	"github.com/yufenghui/go/gravitee/models"
)

var (
	// ErrClientNotFound ...
	ErrClientNotFound = errors.New("Client not found")
	// ErrInvalidClientSecret ...
	ErrInvalidClientSecret = errors.New("Invalid client secret")
	// ErrClientIDTaken ...
	ErrClientIDTaken = errors.New("Client ID taken")
)

//
func (s *Service) FindClientByClientID(clientID string) (*models.OauthClient, error) {
	client := new(models.OauthClient)
	notFound := s.db.Where("client_key = LOWER(?)", clientID).First(client).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrClientNotFound
	}

	return client, nil
}
