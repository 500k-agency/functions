package waitlist

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventID   string          `json:"eventId"`
	EventType string          `json:"eventType"`
	CreatedAt *time.Time      `json:"createdAt"`
	Data      json.RawMessage `json:"data"`
}

type FormResponse struct {
	ResponseID   string       `json:"responseID"`
	SubmissionID string       `json:"submissionID"`
	RespondentID string       `json:"respondentID"`
	FormID       string       `json:"formID"`
	FormName     string       `json:"formName"`
	CreatedAt    *time.Time   `json:"createdAt"`
	Fields       []*FormField `json:"fields"`
}

type FormField struct {
	Key   string      `json:"key"`
	Label string      `json:"label"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	// https://tally.so/help/webhooks
}
