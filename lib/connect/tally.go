package connect

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/500k-agency/function/waitlist"
)

// Tally config struct with exposed methods needed
type Tally struct {
	config Config
}

var (
	TallyClient *Tally

	ErrNotSigned        = errors.New("webhook has no Stripe-Signature header")
	ErrNoValidSignature = errors.New("webhook had no valid signature")
)

// SetupTally sets up stripe with the credentials given
func SetupTally(conf Config) *Tally {
	TallyClient = &Tally{
		config: conf,
	}
	return TallyClient
}

func ComputeSignature(payload []byte, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return mac.Sum(nil)
}

// ConstructEvent validates stripe webhook secret is authentic
func (s *Tally) ConstructEvent(body []byte, header string) (*waitlist.Event, error) {
	t := &waitlist.Event{}

	if header == "" {
		return t, ErrNotSigned
	}

	expectedSignature := base64.StdEncoding.EncodeToString(ComputeSignature(body, s.config.WebhookSecret))
	if header != expectedSignature {
		return t, ErrNoValidSignature
	}

	if err := json.Unmarshal(body, &t); err != nil {
		return t, fmt.Errorf("Failed to parse webhook body json: %s", err.Error())
	}

	return t, nil
}
