package sendgrid

import (
	"context"
	"net/http"
)

type ContactService service

type Contact struct {
	Email               string                 `json:"email"`
	AlternateEmails     []string               `json:"alternate_emails,omitempty"`
	AddressLine1        string                 `json:"address_line_1,omitempty"`
	AddressLine2        string                 `json:"address_line_2,omitempty"`
	City                string                 `json:"city,omitempty"`
	Country             string                 `json:"country,omitempty"`
	FirstName           string                 `json:"first_name,omitempty"`
	LastName            string                 `json:"last_name,omitempty"`
	PostalCode          string                 `json:"postal_code,omitempty"`
	StateProvinceRegion string                 `json:"state_province_region,omitempty"`
	CustomFields        map[string]interface{} `json:"custom_fields,omitempty"`
}

type ContactRequest struct {
	ListIDs  []string   `json:"list_ids"`
	Contacts []*Contact `json:"contacts"`
}

func (s *ContactService) Upsert(ctx context.Context, contactReq *ContactRequest) (JobID, *http.Response, error) {
	req, err := s.client.NewRequest("PUT", "marketing/contacts", contactReq)
	if err != nil {
		return "", nil, err
	}

	jobResponse := map[string]JobID{}
	resp, err := s.client.Do(ctx, req, &jobResponse)
	if err != nil {
		return "", resp, err
	}
	return jobResponse["job_id"], resp, nil
}
