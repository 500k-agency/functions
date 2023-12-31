package sendgrid

import (
	"context"
	"net/http"
)

type MailService service

type MailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Documentation: https://sendgrid.com/docs/api-reference/
type MailPerson struct {
	To                  []*MailAddress         `json:"to"`
	DynamicTemplateData map[string]interface{} `json:"dynamic_template_data"`
	Subject             string                 `json:"subject,omitempty"`

	// Values that are specific to this personalization that will be carried
	// along with the email and its activity data.
	// Substitutions will not be made on custom arguments, so any string that
	// is entered into this parameter will be assumed to be the custom argument
	// that you would like to be used. This field may not exceed 10,000 bytes.
	CustomArgs map[string]string `json:"custom_args"`
}

type MailRequest struct {
	Personalizations []*MailPerson     `json:"personalizations"`
	From             MailAddress       `json:"from"`
	ReplyTo          MailAddress       `json:"reply_to"`
	TemplateID       string            `json:"template_id"`
	Asm              *Asm              `json:"asm,omitempty"`
	MailSettings     *MailSettings     `json:"mail_settings,omitempty"`
	TrackingSettings *TrackingSettings `json:"tracking_settings,omitempty"`
}

// MailSettings defines mail and spamCheck settings
type MailSettings struct {
	SandboxMode *Setting `json:"sandbox_mode,omitempty"`
}

type Asm struct {
	GroupID         int64   `json:"group_id,omitempty"`
	GroupsToDisplay []int64 `json:"groups_to_display,omitempty"`
}

// SubscriptionTrackingSetting ...
type SubscriptionTrackingSetting struct {
	Enable          *bool  `json:"enable,omitempty"`
	Text            string `json:"text,omitempty"`
	Html            string `json:"html,omitempty"`
	SubstitutionTag string `json:"substitution_tag,omitempty"`
}

// TrackingSettings are used to determine how you would like to track the metrics of how your recipients interact with your email.
type TrackingSettings struct {
	SubscriptionTracking *SubscriptionTrackingSetting `json:"subscription_tracking,omitempty"`
}

// Setting enables the mail settings
type Setting struct {
	Enable *bool `json:"enable,omitempty"`
}

// NewSetting ...
func NewSetting(enable bool) *Setting {
	setEnable := enable
	return &Setting{Enable: &setEnable}
}

func (s *MailService) Send(ctx context.Context, mailReq *MailRequest) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "mail/send", mailReq)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
