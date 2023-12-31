package sendgrid

// Documentation: https://sendgrid.com/docs/for-developers/tracking-events/event/
type Event string

const (
	// Delivery events
	EventProcessed = "processed"
	EventDropped   = "dropped"
	EventDelivered = "delivered"
	EventDeferred  = "deferred"
	EventBounce    = "bounce"
	EventBlocked   = "blocked"

	// Engagement events
	EventOpen             = "open"
	EventClick            = "click"
	EventSpamReport       = "spamreport"
	EventUnsubscribe      = "unsubscribe"
	EventGroupUnsubscribe = "group_unsubscribe"
	EventGroupResubscribe = "group_resubscribe"
)

type WebhookEvent struct {
	Attempt           string      `json:"attempt"`
	Category          interface{} `json:"category"`
	Email             string      `json:"email"`
	Event             Event       `json:"event"`
	IP                string      `json:"ip"`
	McStats           string      `json:"mc_stats"`
	PhaseID           string      `json:"phase_id"`
	Response          string      `json:"response"`
	SgContentType     string      `json:"sg_content_type"`
	SgEventID         string      `json:"sg_event_id"`
	SgMessageID       string      `json:"sg_message_id"`
	SgTemplateID      string      `json:"sg_template_id"`
	SgTemplateName    string      `json:"sg_template_name"`
	SinglesendID      string      `json:"singlesend_id"`
	SinglesendName    string      `json:"singlesend_name"`
	SmtpID            string      `json:"smtp-id"`
	TemplateHash      string      `json:"template_hash"`
	TemplateID        string      `json:"template_id"`
	TemplateVersionID string      `json:"template_version_id"`
	Timestamp         int64       `json:"timestamp"`
	Useragent         string      `json:"useragent"`
}
