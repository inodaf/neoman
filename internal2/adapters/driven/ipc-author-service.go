package driven

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/inodaf/neoman/internal2/models"
)

type IPCAuthorService struct{}

type responseDTO struct {
	Name string `json:"name"`
}

func (IPCAuthorService) IsTrusted(authorName string) bool {
  query := url.Values{}
	query.Add("name", url.QueryEscape(authorName))

	resource := url.URL{
		Scheme: "http",
		Host:   "unix",
		Path:   "/author/trust",
    RawQuery: query.Encode(),
	}

	resp, err := UnixSockClient.Get(resource.String())
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func (IPCAuthorService) Trust(authorName string) error {
	resource := url.URL{
		Scheme: "http",
		Host:   "unix",
		Path:   "/author/trust",
	}

	resp, err := UnixSockClient.Post(resource.String(), "text/plain", strings.NewReader(authorName))
	if err != nil {
		return errors.New("could not trust author")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed to trust author")
	}

	return nil
}

func (IPCAuthorService) FindOrCreate(name string) (*models.Author, error) {
	resource := url.URL{
		Scheme: "http",
		Host:   "unix",
		Path:   "/author",
	}

	resp, err := UnixSockClient.Post(resource.String(), "text/plain", strings.NewReader(name))
	if err != nil {
		return nil, errors.New("could not find or create author")
	}

	var authorResponse responseDTO
	err = json.NewDecoder(resp.Body).Decode(&authorResponse)
	if err != nil {
		return nil, err
	}

	return &models.Author{Name: authorResponse.Name}, nil
}
