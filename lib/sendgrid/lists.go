package sendgrid

import (
	"context"
	"net/http"
)

type ListService service

type List struct {
	Name         string `json:"name,omitempty"`
	ID           string `json:"id,omitempty"`
	ContactCount int    `json:"contact_count,omitempty"`
}

func (s *ListService) Create(ctx context.Context, name string) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "marketing/lists", List{Name: name})
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
